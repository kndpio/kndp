const { useState } = window.React;

const Settings: React.FC = () => {
  const [selectedColor, setSelectedColor] = useState('#ababab');
  const [secondColor, setSecondColor] = useState('#242e47');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [logo, setLogo] = useState<string>('');

  const handleColorChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedColor(event.target.value);
  };

  const handleSecondColorChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setSecondColor(event.target.value);
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      const file = event.target.files[0];
      setSelectedFile(file);

      const reader = new FileReader();
      reader.onloadend = () => {
        if (typeof reader.result === 'string') {
          setLogo(reader.result);
        }
      };
      reader.readAsDataURL(file);
    }
  };

  const handleClick = () => {
    setSelectedColor('#ababab');
    setSecondColor('#242e47');
    setSelectedFile(null);

    const data = {
      logo,
      main_color: selectedColor,
      background_color: secondColor,
    };

    console.log(data);

    fetch('http://localhost:3000/profile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.error('Error sending data to the server:', error);
      });
  };
  const handleReload = () => {
    handleClick();
    location.reload();
  };
  return (
    <div
      style={{
        boxSizing: 'border-box',
        position: 'absolute',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        width: '500px',
        height: '700px',
        top: '22%',
        left: '43%',
        borderRadius: '4px',
        border: '1px solid #e8e8e8',
        background: '#fff',
        boxShadow: '1px 2px 3px rgba(0,0,0,.2)',
      }}
    >
      <p
        style={{
          position: 'absolute',
          top: '85px',
          left: '50px',
          fontSize: '20px',
        }}
      >
        Main Color
      </p>
      <input
        style={{
          position: 'absolute',
          marginTop: '70px',
          width: '50px',
          height: '50px',
        }}
        type="color"
        value={selectedColor}
        onChange={handleColorChange}
      />
      <p
        style={{
          position: 'absolute',
          top: '225px',
          left: '50px',
          fontSize: '20px',
        }}
      >
        Background Color
      </p>
      <input
        style={{
          position: 'absolute',
          marginTop: '210px',
          width: '50px',
          height: '50px',
        }}
        type="color"
        value={secondColor}
        onChange={handleSecondColorChange}
      />

      <div>
        <input
          style={{ position: 'absolute', marginTop: '400px', left: '150px' }}
          type="file"
          accept="image/*"
          onChange={handleFileChange}
        />
        {selectedFile && (
          <div style={{ top: '450px', left: '140px', position: 'absolute' }}>
            <h4>Selected Image:</h4>
            <img
              src={URL.createObjectURL(selectedFile)}
              alt="Selected"
              style={{ maxWidth: '250px', height: '143px' }}
            />
          </div>
        )}
      </div>
      <button
        onClick={handleReload}
        style={{
          position: 'absolute',
          top: '640px',
          backgroundColor: '#6d7f8b',
          color: '#FFFFFF',
          padding: '10px 20px',
          border: 'none',
          borderRadius: '5px',
          cursor: 'pointer',
          boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
          fontSize: '16px',
          fontWeight: 'bold',
        }}
      >
        Submit
      </button>
    </div>
  );
};

((window: any) => {
  window?.extensionsAPI?.registerSystemLevelExtension(
    Settings,
    'Profile',
    '/profile',
    'fa-sliders'
  );
})(window);
