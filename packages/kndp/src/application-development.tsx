interface Service {
  name: string;
  parameters: Record<string, any>;
}

const styles = {
  tableList: {
    width: '100%',
    borderCollapse: 'collapse' as 'collapse',
  },
  tableHeader: {
    backgroundColor: '#f0f0f0',
    fontWeight: 'bold',
  },
  actions: {
    backgroundColor: '#f0f0f0',
    fontWeight: 'bold',
  },
  tableCell: {
    padding: '10px',
    border: '1px solid #ccc',
  },
  button: {
    padding: '8px 12px',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    marginRight: '5px',
  },
  submitButton: {
    backgroundColor: '#1fbdd0',
  },
  cards: {
    display: 'flex',
    padding: '-20px',
  },
  cardStyles: {
    backgroundColor: '#f0f0f0',
    borderRadius: '8px',
    padding: '10px',
    margin: '10px',
    width: 'auto',
  },
  name: {
    color: '#1fbdd0',
    fontSize: '25px',
  },
  cardContext: {
    // whiteSpace: 'pre-wrap',
    fontSize: '18px',
    fontFamily: 'Heebo',
  },
};

function formatParameters(parameters: Record<string, any>): string {
  const formattedParameters = Object.keys(parameters).map((key) => {
    const value = parameters[key];
    if (typeof value === 'object') {
      const nestedValue = JSON.stringify(value, null, 2).trim();
      return `${key}:\n${nestedValue.replace(/^{/, '').replace(/}$/, '')}`;
    } else {
      return `${key}:\n${JSON.stringify(value)}`;
    }
  });

  return formattedParameters.join('\n');
}



const TableList: React.FC<{
  services: Service[];
  chartField: string | null;
  version: string | null;
}> = ({ services, chartField, version }) => {
  const handleSubmit = async (
    service: Service,
    chart: string | null,
    version: string | null
  ) => {
    try {
      if (chart && version) {
        const response = await fetch(
          `http://kndp-development-view:8080/api/update/${chart}-${version}`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(service),
          }
        );

        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        console.log(data);
      } else {
        console.error('Chart or version is missing.');
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  return (
    <div style={styles.cards}>
      {services.map((service, index) => (
        <div key={index} style={styles.cardStyles}>
          <h2 style={styles.name}>Name: {service.name}</h2>
          <p>Parameters:</p>
          <pre style={styles.cardContext}>
            {formatParameters(service.parameters)}
          </pre>
          <button
            style={{ ...styles.button, ...styles.submitButton }}
            onClick={() => handleSubmit(service, chartField, version)}
          >
            Submit
          </button>
        </div>
      ))}
    </div>
  );
};

function fetchApplicationDataWithTimeout(timeout: number): Promise<any> {
  return new Promise((resolve) => {
    setTimeout(() => {
      const applicationNameElement = document.querySelector(
        '.argo-dropdown__anchor span'
      );

      if (applicationNameElement) {
        const applicationName = applicationNameElement.textContent.trim();
        const apiUrl = `http://kndp.local/argo-cd/api/v1/applications/${applicationName}`;

        fetch(apiUrl)
          .then((response) => response.json())
          .then((data) => {
            resolve(data);
          })
          .catch((error) => {
            console.error('Error fetching application data:', error);
            // resolve(null);
          });
      } else {
        console.error('Application name element not found.');
        // resolve(null);
      }
    }, timeout);
  });
}

// The main component
const AppDevelopment = () => {
  const [data, setData] = window.React.useState<Service[]>([]);
  const [chartField, setChartField] = window.React.useState<string | null>(
    null
  );
  const [version, setVersion] = window.React.useState<string | null>(null);
  console.log(chartField);
  console.log(version);
  window.React.useEffect(() => {
    fetchApplicationDataWithTimeout(1000).then((data: any) => {
      if (data) {
        const chart = data.spec?.source?.chart;
        setChartField(chart);
        if (data) {
          const targetRevision = data.spec?.source?.targetRevision;
          setVersion(targetRevision);

          if (chart) {
            fetch(
              `http://kndp-development-view:8080/api/stacks/${chart}-${targetRevision}`
            )
              .then((response) => response.json())
              .then((data) => setData(data.services))
              .catch((error) => console.error('Error fetching data:', error));
          }
        }
      }
    });
  }, []);

  return (
    <div>
      <p>Application Development View</p>
      <TableList services={data} chartField={chartField} version={version} />
    </div>
  );
};

((window: any) => {
  window?.extensionsAPI?.registerAppViewExtension(
    AppDevelopment,
    'Development',
    'fa-cogs'
  );
})(window);
