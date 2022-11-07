import React from "react";
import axios from "axios";
import PolicyCreator from "./PolicyCreator";
import PolicySelector from "./PolicySelector";
import Editor from "./Editor";

export default function App() {
  const [policyIds, setPolicyIds] = React.useState([""]);
  const [selectedPolicy, setSelectedPolicy] = React.useState("");
  const [newPolicyAttributes, setNewPolicyAttributes] = React.useState({
    policyKind: "resource",
    name: "",
    version: "",
    scope: "",
  });

  React.useEffect(() => {
    axios
      .get("http://localhost:8090/policies")
      .then((resp) => {
        const ids = resp.data.policyIds;
        if (ids === null) return;
        setPolicyIds(ids);
        if (ids.length && selectedPolicy === "") {
          setSelectedPolicy(ids[0]);
        }
      })
      .catch((e) => {
        console.log(e);
      });
  }, []);

  const handleCreateChange = (e: any) => {
    const key = e.target.id.split("-")[1];
    setNewPolicyAttributes((prevState) => ({
      ...prevState,
      [key]: e.target.value,
    }));
  };

  const updateDisplayedPolicy = (resp: any) => {
    setPolicyIds((prevState) => [...prevState, resp.data.id]);
    setSelectedPolicy(resp.data.id);
  };

  const handleCreatePolicy = (e: any) => {
    e.preventDefault();
    axios
      .post("http://localhost:8090/policy", newPolicyAttributes)
      .then((resp) => {
        updateDisplayedPolicy(resp);
      })
      .catch((e) => {
        console.log(e);
      });
  };

  const handleUpdatePolicy = (policyString: string) => {
    axios
      .patch("http://localhost:8090/policy", { policy: policyString })
      .then((resp) => {
        updateDisplayedPolicy(resp);
      })
      .catch((e) => {
        console.log(e);
      });
  };

  return (
    <div
      id="app"
      className="flex flex-col"
      style={{ paddingLeft: "1.25rem", paddingRight: "1.25rem" }}
    >
      <h1 className="font-medium leading-tight text-4xl mt-0 mb-2">
        Admin API Demo
      </h1>
      <div className="flex flex-col justify-between items-start pb-5">
        <PolicyCreator
          newPolicyAttributes={newPolicyAttributes}
          handleCreateChangeFn={handleCreateChange}
          handleCreatePolicyFn={handleCreatePolicy}
        />
        {policyIds.length ? (
          <PolicySelector
            policyIds={policyIds}
            selectedPolicy={selectedPolicy}
            setSelectedPolicyFn={setSelectedPolicy}
          />
        ) : null}
      </div>
      <Editor
        policyId={selectedPolicy}
        handleUpdatePolicyFn={handleUpdatePolicy}
      />
    </div>
  );
}
