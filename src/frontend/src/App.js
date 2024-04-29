import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Body from "./components/Body/Body";
import Navbar from "./components/Navigation/Nav";

function App() {
  return (
    <div className="dark:bg-slate-900 h-screen ">
      <Navbar />

      <img
        src="images/solvernet_logo.jpeg"
        alt="Sovernet Logo"
        className="logo-background"
      />
      <Body />
      <ToastContainer autoClose={3000} />
    </div>
  );
}

export default App;
