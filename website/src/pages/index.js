import useBaseUrl from '@docusaurus/useBaseUrl';
import React from 'react';

function Home() {
  return (
    <div className="container">
      <header className="hero">
        <img src={useBaseUrl('logo-transparent-small.png')} alt="KNDP Logo" />
        <h1 className="hero__title">Kubernetes Native Development Platform (KNDP)</h1>
        <p className="hero__subtitle">Streamline your Kubernetes-native development.</p>
      </header>

      <main>
        <section>
          <h2>Features</h2>
          <ul>
            <li>Kubernetes-Native Integration</li>
            <li>Developer-Friendly Workflows</li>
            <li>Scalable Architecture</li>
            <li>Extensible Plugin System</li>
          </ul>
        </section>

        <section>
          <h2>Getting Started</h2>
          <p>Learn how to set up and start using KNDP with our comprehensive guides and tutorials.</p>
          <a href={useBaseUrl('docs/getting-started')} className="button">Get Started</a>
        </section>

        <section>
          <h2>Community</h2>
          <p>Join our vibrant community of developers and Kubernetes enthusiasts. Share ideas, ask questions, and collaborate on new features.</p>
          <a href={useBaseUrl('docs/community')} className="button">Join the Community</a>
        </section>
      </main>

      <footer>
        <p>© 2023 KNDP. All rights reserved.</p>
      </footer>
    </div>
  );
}

export default Home;
