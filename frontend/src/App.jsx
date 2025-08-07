import React, { useState } from 'react';
import axios from 'axios';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

const initialYaml = `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-labels
  annotations:
    policies.kyverno.io/title: Require Labels
    policies.kyverno.io/category: Best Practices
    policies.kyverno.io/severity: medium
    policies.kyverno.io/subject: Pod
    policies.kyverno.io/description: >-
      This policy requires that all Pods have the label 'app.kubernetes.io/name'.
spec:
  validationFailureAction: Audit
  background: true
  rules:
    - name: check-for-labels
      match:
        any:
        - resources:
            kinds:
              - Pod
      validate:
        message: "The label 'app.kubernetes.io/name' is required."
        pattern:
          metadata:
            labels:
              app.kubernetes.io/name: "?*"
`;

function App() {
  const [inputYaml, setInputYaml] = useState(initialYaml);
  const [outputYaml, setOutputYaml] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleConvert = async () => {
    setIsLoading(true);
    setError('');
    setOutputYaml('');
    try {
      const response = await axios.post('http://localhost:8080/api/convert', {
        yaml: inputYaml,
      });
      setOutputYaml(response.data.convertedYaml);
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'An unexpected error occurred. Check the backend console.';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 text-gray-200 p-4 sm:p-6 lg:p-8">
      <div className="max-w-7xl mx-auto">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white">Kyverno Policy Converter</h1>
          <p className="mt-2 text-lg text-gray-400">
            Dynamically convert any legacy Kyverno `pattern` or `foreach` policy to a CEL-based `ValidatingPolicy`.
          </p>
        </header>

        <div className="mb-6">
          <button
            onClick={handleConvert}
            disabled={isLoading}
            className="w-full bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-900 disabled:cursor-not-allowed text-white font-bold py-3 px-4 rounded-lg transition-colors duration-300 text-lg"
          >
            {isLoading ? 'Converting...' : 'Transform to CEL'}
          </button>
        </div>

        {error && (
          <div className="bg-red-800 border border-red-600 text-red-100 px-4 py-3 rounded-lg relative mb-6" role="alert">
            <strong className="font-bold">Conversion Error: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}

        <main className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Input Panel */}
          <div className="flex flex-col">
            <h2 className="text-xl font-semibold mb-2 text-white">Legacy Kyverno Policy (YAML Input)</h2>
            <div className="flex-grow bg-gray-800 rounded-lg p-1 border border-gray-700">
              <textarea
                value={inputYaml}
                onChange={(e) => setInputYaml(e.target.value)}
                className="w-full h-96 lg:h-full bg-gray-800 text-gray-200 p-4 rounded-md resize-none font-mono text-sm focus:outline-none"
                placeholder="Paste your legacy Kyverno policy here..."
              />
            </div>
          </div>

          {/* Output Panel */}
          <div className="flex flex-col">
            <h2 className="text-xl font-semibold mb-2 text-white">ValidatingPolicy (CEL Output)</h2>
            <div className="flex-grow bg-[#1e1e1e] rounded-lg overflow-hidden border border-gray-700">
              <SyntaxHighlighter
                language="yaml"
                style={vscDarkPlus}
                customStyle={{ height: '100%', width: '100%', backgroundColor: '#1e1e1e' }}
                codeTagProps={{ style: { fontFamily: "monospace" } }}
              >
                {outputYaml || '# Output will appear here after conversion...'}
              </SyntaxHighlighter>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}

export default App;