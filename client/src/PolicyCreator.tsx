import React from "react";

export default function PolicyCreator(props: {
  newPolicyAttributes: any;
  handleCreateChangeFn: any;
  handleCreatePolicyFn: any;
}) {
  const [showCreator, setShowCreator] = React.useState(false);

  const policyKinds = ["resource", "principal", "derivedRole"];
  const newPolicyFields = ["name", "version", "scope"];

  const handleCreate = (e: any) => {
    props.handleCreatePolicyFn(e);
    setShowCreator((prevState) => !prevState);
  };

  return (
    <div className="flex flex-col mb-2">
      <button
        className="border-2 border-grey hover:border-black font-bold p-1 mb-2 rounded"
        onClick={() => {
          setShowCreator((prevState) => !prevState);
        }}
      >
        Create new policy
      </button>
      <div
        className="flex-col"
        style={{ display: showCreator ? "flex" : "none" }}
      >
        <label>Type:</label>
        <select
          id="newpolicy-policyKind"
          className="p-1"
          value={props.newPolicyAttributes.policyKind}
          onChange={props.handleCreateChangeFn}
        >
          {policyKinds.map((p) => (
            <option key={p} value={p}>
              {p}
            </option>
          ))}
        </select>
        {newPolicyFields.map((f) => {
          return [
            <label key={"label-" + f}>
              {f.charAt(0).toUpperCase() + f.slice(1) + ":"}
            </label>,
            <input
              key={"input-" + f}
              type="text"
              id={"newpolicy-" + f}
              value={props.newPolicyAttributes[f]}
              onChange={props.handleCreateChangeFn}
              className="pl-1 border-solid border-2 border-black rounded"
            />,
          ];
        })}
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold mt-5 mb-5 py-1 px-2 rounded"
          onClick={handleCreate}
        >
          Create
        </button>
      </div>
    </div>
  );
}
