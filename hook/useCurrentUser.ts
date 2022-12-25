import { useRouter } from "next/router";
import { useState, useEffect } from "react";

const useCurrentUser = () => {
  const [currentUser, setCurrentUser] = useState<any>({});
  const router = useRouter();

  useEffect(() => {
    const pathname = window.location.pathname;
    const authToken = localStorage.getItem("talkbox");
    if (pathname === "/") {
      if (authToken) {
        if (router.isReady) {
          const tokens = authToken.split(".");
          setCurrentUser(
            JSON.parse(Buffer.from(tokens[1], "base64").toString())
          );
          router.push("/inboxes");
        }
      }
    } else {
      if (authToken) {
        const tokens = authToken.split(".");
        setCurrentUser(JSON.parse(Buffer.from(tokens[1], "base64").toString()));
      } else {
        if (router.isReady) {
          router.push("/");
        }
      }
    }
  }, [router, router.isReady]);

  return [currentUser, setCurrentUser];
};

export default useCurrentUser;
