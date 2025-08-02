import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/AuthStore";

const Dashboard = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
    }
  }, [isAuthenticated, navigate]);

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Welcome to your dashboard!</h1>
    </div>
  );
};

export default Dashboard;
