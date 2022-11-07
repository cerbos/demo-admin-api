export default function PolicySelector(props: {
  policyIds: any;
  selectedPolicy: string;
  setSelectedPolicyFn: any;
}) {
  return (
    <div className="flex flex-col items-start border-solid border-2 border-grey rounded p-2">
      <label className="pr-2 pb-2 text-xl" htmlFor="policylist">
        Select policy:
      </label>
      <select
        id="policylist"
        className="p-1 w-full"
        value={props.selectedPolicy}
        onChange={(e) => props.setSelectedPolicyFn(e.target.value)}
      >
        {props.policyIds.map((p) => (
          <option key={p} value={p}>
            {p}
          </option>
        ))}
      </select>
    </div>
  );
}
