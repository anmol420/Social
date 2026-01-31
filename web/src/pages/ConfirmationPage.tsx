import { useNavigate, useParams } from "react-router-dom";
import { API_URL } from "../App";

function ConfirmationPage() {
  const { token } = useParams();
  const redirect = useNavigate();
  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT",
    });
    if (response.ok) {
      redirect("/");
    } else {
      alert("Failed to confirm token");
    }
  };

  return (
    <>
      <h1>Confirmation</h1>
      <button onClick={handleConfirm}>Click To Confirm</button>
    </>
  );
}

export default ConfirmationPage;
