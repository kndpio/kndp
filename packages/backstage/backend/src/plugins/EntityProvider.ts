import { UrlReader } from '@backstage/backend-common';
import { Entity } from '@backstage/catalog-model';
import { EntityProvider, EntityProviderConnection } from '@backstage/plugin-catalog-node';
import { readFileSync } from 'fs';

export class KubernetesProvider implements EntityProvider {
  private readonly env: string;
  private readonly reader: UrlReader;
  private connection?: EntityProviderConnection;

  private readonly apiEndpoints: string[]; // Define API endpoints here

  constructor(env: string, reader: UrlReader) {
    this.env = env;
    this.reader = reader;
    this.apiEndpoints = (process.env.API_ENDPOINTS || '').split(';').filter(Boolean);
  }
  
  getProviderName(): string {
    return `kubernetes-${this.env}`;
  }
  async connect(connection: EntityProviderConnection): Promise<void> {
    this.connection = connection;
  }
 
  async run(): Promise<void> {
    if (!this.connection) {
      throw new Error('Not initialized');
    }

    try {
      const serviceAccountToken = readFileSync('/var/run/secrets/kubernetes.io/serviceaccount/token', 'utf8');

      for (const apiUrl of this.apiEndpoints) {
        try {
          const fetchOptions = {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${serviceAccountToken}`
            }
          };
          const request = new Request(apiUrl, fetchOptions);
          console.log(`Making API request to: ${apiUrl}`);

          const response = await fetch(request);

          if (response.ok) {
            const data = await response.text();
            const entities = processKubernetesData(data);
            await this.connection.applyMutation({
              type: 'full',
              entities: entities.map(entity => ({
                entity,
                locationKey: `kubernetes-provider:${this.env}`,
              })),
            });
          } else {
            console.error('Request failed with status: ' + response.status);
          }
        } catch (error) {
          console.error('Error making API request:', error);
        }
      }
    } catch (error) {
      console.error('Error:', error);
    }
  }
}

function processKubernetesData(data: string): Entity[] {
  try {
    const entities: Entity[] = [];
    const parsedData = JSON.parse(data);
    if (parsedData) {
      const entity = {
        kind: 'Component',
        apiVersion: 'backstage.io/v1alpha1',
        metadata: {
          annotations: {
            "backstage.io/managed-by-location": "url:https://kubernetes.default.svc/apis/kndp.io/v1alpha1/releases",
            "backstage.io/managed-by-origin-location": "url:https://kubernetes.default.svc/apis/kndp.io/v1alpha1/releases"
          },
          name: parsedData.metadata.name,
          namespace: parsedData.metadata.namespace,
        },
        spec: {
          type: parsedData.kind,
          lifecycle: 'experimental',
          owner: 'guests'
        },
      };

      entities.push(entity);
    }
    return entities;
  } catch (error) {
    console.error('Error processing Kubernetes data:', error);
    return [];
  }
}
