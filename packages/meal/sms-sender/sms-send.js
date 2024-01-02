import { CustomObjectsApi, KubeConfig } from '@kubernetes/client-node';
import fs from 'fs';
import RouterClient from './src/routerClient.mjs';
import { TP_ACT, TP_CONTROLLERS } from './src/routerProtocol.mjs';
import dotenv from 'dotenv';
dotenv.config();

const group = 'kndp.io';
const version = 'v1alpha1';
const plural = 'meals';

const kubeConfig = new KubeConfig();
kubeConfig.loadFromDefault();
const client = kubeConfig.makeApiClient(CustomObjectsApi);

// getting meal object from k8s api
const resourceName = process.env.MEAL_NAME || 'meal';
let response = {};

try {
  response = await client.getClusterCustomObject(
    group,
    version,
    plural,
    resourceName
  );
} catch (error) {
  console.error(
    'Error updating custom resource:',
    error.response?.body || error.message
  );
}

// getting number of employees with 'yes' status
let employee_refs = response.body.spec.employeeRefs;

function countUsers() {
  let count = 0;
  employee_refs.forEach((employee_ref) => {
    if (employee_ref.status.toLowerCase() === 'yes') {
      count++;
    }
  });
  return count;
}

const textContent = countUsers();
console.log(textContent);

let orderTime = response.body.spec.dueOrderTime;

// Sending the message to restaurant with total number of meals if orderTime is 'over'
if (orderTime === 'over') {
  console.log('Sms is ready to be sent:', orderTime);

  let configFilePath = 'config.json';
  let rawConfig = fs.readFileSync(configFilePath);
  let config = JSON.parse(rawConfig);
  let routerUiUrl = process.env.URL || config.url;
  let routerUiLogin = process.env.LOGIN || config.login;
  let routerUiPassword = process.env.PASSWORD || config.password;
  const to = process.env.TO || config.to;
  console.log(' router-url: ', routerUiUrl, '\n' , 'router-login: ', routerUiLogin ,'\n', 'router-password: ', routerUiPassword, '\n' , 'destination-number: ', to);

  try {
    routerUiUrl;
    routerUiLogin;
    routerUiPassword;
  } catch (exception) {
    console.log(
      'config file ' + configFilePath + ' could not be read, skipping'
    );
  }

  const client_router = new RouterClient(
    routerUiUrl,
    routerUiLogin,
    routerUiPassword
  );

  const payloadSendSms = {
    method: TP_ACT.ACT_SET,
    controller: TP_CONTROLLERS.LTE_SMS_SENDNEWMSG,
    attrs: {
      index: 1,
      to,
      textContent,
    },
  };

  const payloadGetSendSmsResult = {
    method: TP_ACT.ACT_GET,
    controller: TP_CONTROLLERS.LTE_SMS_SENDNEWMSG,
    attrs: ['sendResult'],
  };

  client_router
    .connect()
    .then((_) => client_router.execute(payloadSendSms))
    .then(verify_submission)
    .then((_) => client_router.execute(payloadGetSendSmsResult))
    .then(verify_submission_result)
    .then((_) => client_router.disconnect())
    .catch(function (error) {
      // handle error, exit failure
      console.log(error);
      process.exit(1);
    });

  function verify_submission(result) {
    if (result.error === 0) {
      console.log('Great! SMS send operation was accepted.');
    } else {
      // hopefully we will never have this error
      throw new Error('SMS send operation was not accepted');
    }
  }

  function verify_submission_result(result) {
    if (result.error === 0 && result.data[0]['sendResult'] === 1) {
      console.log('Great! SMS sent successfully');
    } else if (result.error === 0 && result.data[0]['sendResult'] === 3) {
      //TODO sendResult=3 means queued or processing ??
      console.log('Warning: SMS sending was accepted but not yet processed.');
    } else {
      console.log('Error: SMS could not be sent by router');
    }
  }
} else {
  console.log('dueOrderTime is not over');
}
