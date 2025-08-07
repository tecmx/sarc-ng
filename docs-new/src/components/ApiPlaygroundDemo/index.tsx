import React, { useState } from 'react';
import styles from './styles.module.css';

/**
 * A simplified API playground component for demonstration purposes
 * In a real implementation, this would actually call the API
 */
export default function ApiPlaygroundDemo(): JSX.Element {
  const [endpoint, setEndpoint] = useState<string>('/api/v1/buildings');
  const [method, setMethod] = useState<string>('GET');
  const [requestBody, setRequestBody] = useState<string>('');
  const [response, setResponse] = useState<{status: number, body: string} | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  // Simulated API responses
  const mockResponses: Record<string, Record<string, {status: number, body: any}>> = {
    '/api/v1/buildings': {
      GET: {
        status: 200,
        body: [
          {
            "id": 1,
            "name": "Main Campus Building",
            "code": "MCB-A",
            "createdAt": "2024-01-01T00:00:00Z",
            "updatedAt": "2024-01-01T00:00:00Z"
          },
          {
            "id": 2,
            "name": "Science Building",
            "code": "SCI-B",
            "createdAt": "2024-01-01T00:00:00Z",
            "updatedAt": "2024-01-01T00:00:00Z"
          }
        ]
      },
      POST: {
        status: 201,
        body: {
          "id": 3,
          "name": "New Building",
          "code": "NEW-C",
          "createdAt": "2024-06-01T00:00:00Z",
          "updatedAt": "2024-06-01T00:00:00Z"
        }
      }
    },
    '/api/v1/resources': {
      GET: {
        status: 200,
        body: [
          {
            "id": 1,
            "name": "Projector Room A",
            "type": "equipment",
            "isAvailable": true,
            "createdAt": "2024-01-01T00:00:00Z",
            "updatedAt": "2024-01-01T00:00:00Z"
          }
        ]
      }
    }
  };

  const handleSend = () => {
    setLoading(true);

    // Simulate API request
    setTimeout(() => {
      const mockResponse = mockResponses[endpoint]?.[method];
      if (mockResponse) {
        setResponse({
          status: mockResponse.status,
          body: JSON.stringify(mockResponse.body, null, 2)
        });
      } else {
        setResponse({
          status: 404,
          body: JSON.stringify({ error: "Endpoint not found in demo" }, null, 2)
        });
      }
      setLoading(false);
    }, 500);
  };

  return (
    <div className={styles.apiPlayground}>
      <h3>API Playground Demo</h3>

      <div className={styles.controls}>
        <select
          value={method}
          onChange={(e) => setMethod(e.target.value)}
          className={styles.methodSelect}
        >
          <option value="GET">GET</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="DELETE">DELETE</option>
        </select>

        <select
          value={endpoint}
          onChange={(e) => setEndpoint(e.target.value)}
          className={styles.endpointSelect}
        >
          <option value="/api/v1/buildings">/api/v1/buildings</option>
          <option value="/api/v1/resources">/api/v1/resources</option>
        </select>

        <button
          onClick={handleSend}
          disabled={loading}
          className={styles.sendButton}
        >
          {loading ? 'Sending...' : 'Send'}
        </button>
      </div>

      {method !== 'GET' && (
        <div className={styles.requestBodyContainer}>
          <div className={styles.label}>Request Body:</div>
          <textarea
            value={requestBody}
            onChange={(e) => setRequestBody(e.target.value)}
            placeholder={`{\n  "name": "New Building",\n  "code": "NEW-C"\n}`}
            className={styles.textarea}
          />
        </div>
      )}

      {response && (
        <div className={styles.responseContainer}>
          <div className={styles.responseHeader}>
            Response: <strong>Status {response.status}</strong>
          </div>
          <pre className={styles.responseBody}>
            {response.body}
          </pre>
        </div>
      )}
    </div>
  );
}
