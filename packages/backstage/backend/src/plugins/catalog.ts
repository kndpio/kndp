import { CatalogBuilder } from '@backstage/plugin-catalog-backend';
import { Router } from 'express';
import { PluginEnvironment } from '../types';
import { KubernetesProvider } from './EntityProvider';

export default async function createPlugin(
  env: PluginEnvironment,
): Promise<Router> {
  const builder = await CatalogBuilder.create(env);
  const kubernetes = new KubernetesProvider('production', env.reader);
  builder.addEntityProvider(kubernetes); // Add the entity provider
  const { processingEngine, router } = await builder.build();
  await processingEngine.start();
  await env.scheduler.scheduleTask({
    id: 'run_kubernetes_refresh',
    fn: async () => {
      await kubernetes.run();
    },
    frequency: { minutes: 30 },
    timeout: { minutes: 10 },
  });
  return router;
}
