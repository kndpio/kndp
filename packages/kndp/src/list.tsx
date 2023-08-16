interface Item {
  id: string;
  name: string;
  status: string;
  description: string;
  isActive: boolean;
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
    width: '220px',
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
  activateButton: {
    backgroundColor: '#1fbdd0',
  },
  deactivateButton: {
    backgroundColor: '#6d7f8b',
  },
};

const TableList: React.FC<{ items: Item[] }> = ({ items }) => {
  const handleActivate = (id: string) => {
    fetch('http://192.168.88.249:8080/api/stacks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id,
        activate: true,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.error('Error activating item:', error);
      });

    console.log(`Activate item with ID: ${id}`);
    setTimeout(() => {
      location.reload();
    }, 1000);
  };

  const handleDeactivate = (id: string) => {
    fetch('http://192.168.88.249:8080/api/stacks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id,
        activate: false,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.error('Error activating item:', error);
      });

    console.log(`Deactivate item with ID: ${id}`);
    setTimeout(() => {
      location.reload();
    }, 1000);
  };
  return (
      <div
        style={{
          boxSizing: 'border-box',
          position: 'absolute',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          width: '70%',
          height: 'auto',
          top: '22%',
          left: '20%',
          borderRadius: '4px',
          border: '1px solid #e8e8e8',
          background: '#fff',
          boxShadow: '1px 2px 3px rgba(0,0,0,.2)',
        }}
      >
        <table style={styles.tableList}>
          <thead>
            <tr>
              <th style={styles.tableHeader}>ID</th>
              <th style={styles.tableHeader}>Name</th>
              <th style={styles.tableHeader}>Status</th>
              <th style={styles.tableHeader}>Description</th>
              <th style={styles.actions}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {items.map((item) => (
              <tr key={item.id}>
                <td style={styles.tableCell}>{item.id}</td>
                <td style={styles.tableCell}>{item.name}</td>
                <td style={styles.tableCell}>{item.status}</td>
                <td style={styles.tableCell}>{item.description}</td>
                <td style={styles.tableCell}>
                  <button
                    style={{ ...styles.button, ...styles.activateButton }}
                    onClick={() => handleActivate(item.id)}
                  >
                    Activate
                  </button>

                  <button
                    style={{ ...styles.button, ...styles.deactivateButton }}
                    onClick={() => handleDeactivate(item.id)}
                  >
                    Deactivate
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
  );
};

const App: React.FC = () => {
  const [data, setData] = window.React.useState<Item[]>([]);

  window.React.useEffect(() => {
    fetch('http://192.168.88.249:8080/api/stacks')
      .then((response) => response.json())
      .then((data) => setData(data))
      .catch((error) => console.error('Error fetching data:', error));
  }, []);

  return (
    <div>
      <h1 style={{ color: 'black', textAlign: 'center' }}>List</h1>
      <TableList items={data} />
    </div>
  );
};

((window: any) => {
  window?.extensionsAPI?.registerSystemLevelExtension(
    App,
    'List',
    '/list',
    'fa-list-ul'
  );
})(window);
