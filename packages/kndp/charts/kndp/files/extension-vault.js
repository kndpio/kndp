((window) => {
    const component = () => {
      const iframeProps = {
        src: "https://kndp.local/ui", 
        width: "100%",
        height: "890px", 
      };
      return React.createElement("iframe", iframeProps);
    };
  
    window.extensionsAPI.registerSystemLevelExtension(
      component,
      "Vault",
      "/vault",
      "fa fa-key "
    );
  
  })(window);
  