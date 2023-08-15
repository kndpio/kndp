((window) => {
  const component = () => {
    const iframeRef = React.useRef(null);

    const adjustIframeHeight = () => {
      if (iframeRef.current) {
        const windowHeight = window.innerHeight;
        iframeRef.current.style.height = windowHeight + 'px';
      }
    };

    React.useEffect(() => {
      adjustIframeHeight(); 
      window.addEventListener('resize', adjustIframeHeight);
      return () => {
        window.removeEventListener('resize', adjustIframeHeight);
      };
    }, []);

    const iframeProps = {
      src: "https://kndp.local/ui",
      width: "100%",
      height: "100%", 
      ref: iframeRef, 
    };
    return React.createElement("iframe", iframeProps);
  };

  window.extensionsAPI.registerSystemLevelExtension(
    component,
    "Secrets",
    "/vault",
    "fa fa-key "
  );
})(window);
