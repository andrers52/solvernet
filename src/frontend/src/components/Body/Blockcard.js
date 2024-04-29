import React, { useState, useEffect } from "react";
import axios from "axios";
import { toast } from "react-toastify";
import loader from "../../assets/loader.gif";
import { apiBaseUrl } from "../../Utils/Apis";

function Blockcard({ reload }) {
  const [data, setdata] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    FetchBlocks();
    // const intervalId = setInterval(() => {
    //   FetchBlocks();
    // }, 3000); // Fetches new blockchain data every 3 seconds

    // // Cleanup function to clear the interval
    // return () => clearInterval(intervalId);
  }, [reload]); // Dependency array ensures the interval is reset when 'reload' changes

  const FetchBlocks = async () => {
    setLoading(true);
    await axios
      .get(apiBaseUrl + "/api/get_blockchain")
      .then((res) => {
        setdata(res?.data);
        setLoading(false);
      })
      .catch((error) => {
        toast.error(`error occured: ${error}`, {
          position: toast.POSITION.TOP_RIGHT,
        });
        setLoading(false);
      });
  };

  return (
    <div className="flex flex-col align-center justify-center p-6">
      <div className="flex items-center justify-center h-5">
        {loading && (
          <p className="dark:text-white">
            <b className="flex items-center justify-center">Loading .....</b>
          </p>
        )}
      </div>

      {!loading && (
        <div className="grid-cols-1 grid sm:grid-cols-4 grid-flow-row-dense gap-4 ">
          {data?.map((bc, i) => {
            let backgroundColor;

            switch (bc.data.type) {
              case 0:
                backgroundColor = "red";
                break;
              case 1:
                backgroundColor = "blue";
                break;
              case 2:
                backgroundColor = "green";
                break;
              default:
                backgroundColor = "black";
            }
            const bannerStyle = {
              backgroundColor: backgroundColor,
              height: "20px", // Adjust height as needed
              borderTopLeftRadius: "5px",
              borderTopRightRadius: "5px",
            };
            return (
              <div
                className="p-4 border-[1px] border-[#3d4f7c] rounded-md hover:bg-slate-50 dark:text-white dark:hover:bg-slate-800"
                key={i}
              >
                <div style={bannerStyle}></div>
                <p>Block ID: {bc.index} </p>
                <p className="break-all">Previous Hash: {bc.prevhash}</p>
                <p className="break-all">Block Hash :{bc.hash}</p>
                {bc.data.type === 0 && (
                  <div>
                    <p>Transaction</p>
                    <p>Recipient: {bc.data.transaction.to}</p>
                    <p>Sender: {bc.data.transaction.from}</p>
                    <p>Amount: {bc.data.transaction.amount}</p>
                    <p>
                      ProblemBlockHeight:{" "}
                      {bc.data.transaction.problem_block_height}
                    </p>
                  </div>
                )}
                {bc.data.type === 1 && (
                  <div>
                    <p>Problem</p>
                    <p>Items: {JSON.stringify(bc.data.problem.items)}</p>
                    <p>Capacity: {bc.data.problem.capacity}</p>
                    <p>Bounty: {bc.data.problem.bounty}</p>
                    <p>Address: {bc.data.problem.address}</p>
                  </div>
                )}
                {bc.data.type === 2 && (
                  <div>
                    <p>Solution</p>
                    <p>
                      Items: {JSON.stringify(bc.data.proposed_solution.items)}
                    </p>
                    <p>
                      ProblemBlockHeight:{" "}
                      {bc.data.proposed_solution.problem_block_height}
                    </p>
                    <p>Value: {bc.data.proposed_solution.value}</p>
                    <p>Address: {bc.data.proposed_solution.address}</p>
                  </div>
                )}
                <br />
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

export default Blockcard;
