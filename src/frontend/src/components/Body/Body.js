import React, { useState } from "react";
import Blockcard from "./Blockcard";
import loader from "../../assets/loader.gif";
import axios from "axios";
import { toast } from "react-toastify";

import { apiBaseUrl } from "../../Utils/Apis";

function Body() {
  const [stateLoad, setStateLoad] = useState(true);

  const [blockHeight, setBlockHeight] = useState(0);
  const [indexes, setIndexes] = useState("");
  const [value, setValue] = useState(0);
  const [address, setAddress] = useState("");

  const [bounty, setBounty] = useState(0);
  const [problemAddress, setProblemAddress] = useState("");
  const [capacity, setCapacity] = useState(0);
  const [weights, setWeights] = useState("");
  const [values, setValues] = useState("");

  const NewSolution = async () => {
    if (
      blockHeight === "" ||
      indexes === "" ||
      value === "" ||
      address === ""
    ) {
      toast.error(`Please fill out all the fields`, {
        position: toast.POSITION.TOP_RIGHT,
      });
    } else {
      setStateLoad(false);

      const postMe = JSON.stringify({
        block_height: Number(blockHeight),
        items: indexes
          .split(",")
          .map((substring) => parseInt(substring.trim())),
        value: Number(value),
        address: address,
      });
      console.log(postMe);
      await axios
        .post(apiBaseUrl + "/api/send_proposed_solution", postMe)
        .then((res) => {
          toast.info(`Sent successfully`, {
            position: toast.POSITION.TOP_RIGHT,
          });

          setStateLoad(true);
        })
        .catch((error) => {
          toast.error(`error occured: ${error}`, {
            position: toast.POSITION.TOP_RIGHT,
          });
          setStateLoad(false);
        });
    }
  };

  const NewProblem = async () => {
    if (
      bounty === "" ||
      problemAddress === "" ||
      capacity === "" ||
      weights === "" ||
      values === ""
    ) {
      toast.error(`Please fill out all the fields`, {
        position: toast.POSITION.TOP_RIGHT,
      });
    } else {
      setStateLoad(false);
      const parsedWeights = weights.split(",");
      const parsedValues = values.split(",");

      // Initialize an empty array to store items
      const items = [];

      // Ensure the lengths of both arrays match
      if (parsedWeights.length !== parsedValues.length) {
        throw new Error("Number of weights doesn't match number of values.");
      }

      // Iterate over the arrays to create items
      for (let i = 0; i < parsedWeights.length; i++) {
        // Convert string values to numbers
        const weight = parseInt(parsedWeights[i]);
        const value = parseInt(parsedValues[i]);

        // Create an item object with weight and value properties
        const item = {
          weight: weight,
          value: value,
        };

        // Push the item to the items array
        items.push(item);
      }

      // Return the array of items as JSON
      const postMe = JSON.stringify({
        bounty: Number(bounty),
        address: problemAddress,
        capacity: Number(capacity),
        items: items,
      });
      console.log(postMe);

      await axios
        .post(apiBaseUrl + "/api/send_problem", postMe)
        .then((res) => {
          toast.info(`Sent successfully`, {
            position: toast.POSITION.TOP_RIGHT,
          });

          setStateLoad(true);
        })
        .catch((error) => {
          toast.error(`error occured: ${error}`, {
            position: toast.POSITION.TOP_RIGHT,
          });
          setStateLoad(false);
        });
    }
  };



  return (
    <div className="dark:bg-slate-900">
      <div className="flex flex-col flex-1 items-center justify-start w-full mf:mt-0 mt-10">
        <h1 className="text-gray-600 text-3xl mt-2 mb-5 text-gradient dark:text-white">
          SolverNet
        </h1>
        <table border="1">
          <tr>
            <td>
              <div className="p-5 sm:w-96 w-full flex flex-col justify-start items-center blue-glassmorphism bg-slate-50 dark:bg-slate-900 shadow-sm">
                <h1 className="dark:text-white text-sm mt-4">
                  Propose a solution
                </h1>
                <input
                  onChange={(e) => setBlockHeight(e.target.value)}
                  placeholder="Block you want to submit a solution for"
                  type="number"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setIndexes(e.target.value)}
                  placeholder="Item indexes of your solution (comma separated)"
                  type="text"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setValue(e.target.value)}
                  placeholder="Value of your solution"
                  type="number"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setAddress(e.target.value)}
                  placeholder="Your address"
                  type="text"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />

                <div className="h-[1px] w-full bg-gray-400 my-2" />

                <button
                  type="button"
                  onClick={() => NewSolution()}
                  className="dark:text-white w-full mt-2 border-[1px] p-2 border-[#3d4f7c] hover:bg-[#dce4f1] dark:hover:bg-slate-800 rounded-full cursor-pointer"
                >
                  Send
                </button>
              </div>
            </td>
            <td>
              <div className="p-5 sm:w-96 w-full flex flex-col justify-start items-center blue-glassmorphism bg-slate-50 dark:bg-slate-900 shadow-sm">
                <h1 className="dark:text-white text-sm mt-4">
                  Propose a problem
                </h1>

                <input
                  onChange={(e) => setBounty(e.target.value)}
                  placeholder="Bounty"
                  type="number"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setProblemAddress(e.target.value)}
                  placeholder="Your address"
                  type="text"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setCapacity(e.target.value)}
                  placeholder="Knapsack capacity"
                  type="number"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setWeights(e.target.value)}
                  placeholder="Item weights (comma separated)"
                  type="text"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />
                <input
                  onChange={(e) => setValues(e.target.value)}
                  placeholder="Item values (comma separated)"
                  type="text"
                  className="my-2 w-full rounded-md p-2 outline-none bg-transparent dark:text-white border-[1px] border-[#3d4f7c]  text-sm white-glassmorphism"
                />

                <div className="h-[1px] w-full bg-gray-400 my-2" />

                <button
                  type="button"
                  onClick={() => NewProblem()}
                  className="dark:text-white w-full mt-2 border-[1px] p-2 border-[#3d4f7c] hover:bg-[#dce4f1] dark:hover:bg-slate-800 rounded-full cursor-pointer"
                >
                  Send
                </button>
              </div>
            </td>
          </tr>
        </table>
      </div>

      <div className="txresult flex align-center justify-center w-full mb-10 h-2 dark:text-gray-400">
        {!stateLoad && <span>Processing ...</span>}
      </div>

      <Blockcard reload={stateLoad} />
    </div>
  );
}

export default Body;
