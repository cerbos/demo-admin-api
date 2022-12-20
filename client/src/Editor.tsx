// @ts-ignore
import JSONWorker from 'url:monaco-editor/esm/vs/language/json/json.worker.js';
// @ts-ignore
import CSSWorker from 'url:monaco-editor/esm/vs/language/css/css.worker.js';
// @ts-ignore
import HTMLWorker from 'url:monaco-editor/esm/vs/language/html/html.worker.js';
// @ts-ignore
import TSWorker from 'url:monaco-editor/esm/vs/language/typescript/ts.worker.js';
// @ts-ignore
import EditorWorker from 'url:monaco-editor/esm/vs/editor/editor.worker.js';
import React from "react";
import axios from "axios";
import * as monaco from 'monaco-editor';

self.MonacoEnvironment = {
	getWorkerUrl: function (_, label: string) {
		if (label === 'json') {
			return JSONWorker;
		}
		if (label === 'css' || label === 'scss' || label === 'less') {
			return CSSWorker;
		}
		if (label === 'html' || label === 'handlebars' || label === 'razor') {
			return HTMLWorker;
		}
		if (label === 'typescript' || label === 'javascript') {
			return TSWorker;
		}
		return EditorWorker;
	}
};

export default function Editor(props: {
  policyId: string;
  handleUpdatePolicyFn: any;
}) {
  const container = React.useRef(null);
  const editor = React.useRef(null);
  const [isValid, setIsValid] = React.useState(true);
  const [errString, setErrString] = React.useState("");

  const handleTextChange = () => {
    const text = editor.current.getModel().getValue();
    axios
      .patch("http://localhost:8090/validate", { policy: text })
      .then(() => {
        setIsValid(true);
      })
      .catch((e) => {
        setErrString(e.response.data);
        setIsValid(false);
      });
  };

  React.useEffect(() => {
    if (editor.current !== null) return;
    editor.current = monaco.editor.create(container.current, {
      value: "",
      language: "yaml",
    });

    // attach listener
    editor.current.getModel().onDidChangeContent(() => {
      handleTextChange();
    });
  }, []);

  React.useEffect(() => {
    if (editor.current === null) return;
    if (props.policyId === "") return;
    axios
      .get("http://localhost:8090/policy?id=" + props.policyId)
      .then((resp) => {
        editor.current.getModel().setValue(resp.data.policies[0]);
      })
      .catch((e) => {
        console.log(e);
      });
  }, [props.policyId]);

  const handleUpdate = () => {
    const text = editor.current.getModel().getValue();
    props.handleUpdatePolicyFn(text);
  };

  return (
    <div>
      <div
        ref={container}
        style={{ height: "400px", border: "1px solid grey" }}
      ></div>
      <div className="flex flex-row font-medium mt-5 border-solid border-2 border-grey rounded p-2">
        <button
          disabled={!isValid}
          className={
            isValid
              ? "bg-blue-500 hover:bg-blue-700 text-white font-bold px-2 py-1 mr-3 rounded"
              : "bg-gray-300 text-white font-bold px-2 py-1 mr-3 rounded"
          }
          onClick={handleUpdate}
        >
          {isValid ? "Save" : "Invalid"}
        </button>
        {isValid ? (
          <p className="flex items-center">✅ Policy valid</p>
        ) : (
          <p className="flex items-center">❌ Policy invalid: {errString}</p>
        )}
      </div>
    </div>
  );
}
