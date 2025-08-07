import React, { useState } from 'react';
import styles from './styles.module.css';
import clsx from 'clsx';

interface Endpoint {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  path: string;
  description: string;
}

interface ApiEndpointDemoProps {
  endpoints: Endpoint[];
}

export default function ApiEndpointDemo({ endpoints }: ApiEndpointDemoProps): JSX.Element {
  const [activeMethod, setActiveMethod] = useState<string | null>(null);

  // Group endpoints by their primary resource
  const groupedEndpoints = endpoints.reduce((acc, endpoint) => {
    // Extract the resource from the path (e.g., /buildings, /resources)
    const resource = endpoint.path.split('/')[1];
    if (!acc[resource]) {
      acc[resource] = [];
    }
    acc[resource].push(endpoint);
    return acc;
  }, {} as Record<string, Endpoint[]>);

  const getMethodClass = (method: string) => {
    switch (method) {
      case 'GET': return styles.get;
      case 'POST': return styles.post;
      case 'PUT': return styles.put;
      case 'DELETE': return styles.delete;
      default: return '';
    }
  };

  const toggleMethod = (method: string) => {
    if (activeMethod === method) {
      setActiveMethod(null);
    } else {
      setActiveMethod(method);
    }
  };

  const filteredEndpoints = activeMethod
    ? endpoints.filter(e => e.method === activeMethod)
    : endpoints;

  return (
    <div className={styles.container}>
      <div className={styles.methodFilters}>
        <button
          className={clsx(styles.methodButton, styles.get, { [styles.active]: activeMethod === 'GET' })}
          onClick={() => toggleMethod('GET')}
        >
          GET
        </button>
        <button
          className={clsx(styles.methodButton, styles.post, { [styles.active]: activeMethod === 'POST' })}
          onClick={() => toggleMethod('POST')}
        >
          POST
        </button>
        <button
          className={clsx(styles.methodButton, styles.put, { [styles.active]: activeMethod === 'PUT' })}
          onClick={() => toggleMethod('PUT')}
        >
          PUT
        </button>
        <button
          className={clsx(styles.methodButton, styles.delete, { [styles.active]: activeMethod === 'DELETE' })}
          onClick={() => toggleMethod('DELETE')}
        >
          DELETE
        </button>
      </div>

      {Object.keys(groupedEndpoints).map(resource => (
        <div key={resource} className={styles.resourceSection}>
          <h4 className={styles.resourceTitle}>/{resource}</h4>
          <div className={styles.endpointList}>
            {groupedEndpoints[resource]
              .filter(e => !activeMethod || e.method === activeMethod)
              .map((endpoint, i) => (
                <div key={i} className={styles.endpoint}>
                  <span className={clsx(styles.method, getMethodClass(endpoint.method))}>
                    {endpoint.method}
                  </span>
                  <span className={styles.path}>
                    {endpoint.path}
                  </span>
                  <span className={styles.description}>
                    {endpoint.description}
                  </span>
                </div>
              ))}
          </div>
        </div>
      ))}
    </div>
  );
}
