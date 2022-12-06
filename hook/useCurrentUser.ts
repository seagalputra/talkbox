import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";

const useCurrentUser = () => {
  const [cookies] = useCookies(["talkbox"]);
  const [currentUser, setCurrentUser] = useState<any>({});

  useEffect(() => {
    if (cookies) {
      const [, payload] = cookies.talkbox?.split(".");
      setCurrentUser(JSON.parse(Buffer.from(payload, "base64").toString()));
    }
  }, [cookies]);

  return [currentUser, setCurrentUser];
};

export default useCurrentUser;
