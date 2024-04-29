import React, { useState } from "react";
import Blockcard from "./Blockcard";

function Body() {
  const [stateLoad, setStateLoad] = useState(true);

  return (
    <div className="dark:bg-slate-900">
      <div className="flex flex-col flex-1 items-center justify-start w-full mf:mt-0 mt-10">
        <h1 className="text-gray-600 text-3xl mt-2 mb-5 text-gradient dark:text-white">
          SolverNet
        </h1>
      </div>

      <div className="txresult flex align-center justify-center w-full mb-10 h-2 dark:text-gray-400">
        {!stateLoad && <span>Processing ...</span>}
      </div>

      <Blockcard reload={stateLoad} />
    </div>
  );
}

export default Body;
