import Link from "next/link";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import { useEffect, useState, useRef, useContext } from "react";
import type { NextPageWithLayout } from "../_app";
import type { ReactElement } from "react";
import InboxesLayout, { RoomContext } from "./_layout";
import http from "../../lib/http";
import useCurrentUser from "../../hook/useCurrentUser";

type SendMessageInput = {
  body: string;
  attachment?: string;
};

const Inboxes: NextPageWithLayout<any> = () => {
  const router = useRouter();
  const [wsInstance, setWsInstance] = useState<any>(null);
  const { register, handleSubmit, formState, reset } =
    useForm<SendMessageInput>();
  const [currentUser, setCurrentUser] = useCurrentUser();
  const [messages, setMessages] = useState<any>([]);
  const [isFetchingMessage, setIsFetchingMessage] = useState<boolean>(false);
  const messageBoxRef = useRef<any>();
  const { getCurrentRoom } = useContext(RoomContext);

  useEffect(() => {
    if (router.isReady) {
      const { roomId } = router.query;
      const url = `${process.env.NEXT_PUBLIC_WS_BASE_URL}/rooms/${roomId}`;
      let ws: any;
      if (typeof window !== "undefined") {
        ws = new WebSocket(url);
        setWsInstance(ws);
      }

      return () => {
        if (ws?.readyState !== 3) {
          ws.close();
        }
      };
    }
  }, [router.isReady, router.query]);

  useEffect(() => {
    if (formState?.isSubmitSuccessful) {
      reset();
    }
  }, [formState, reset]);

  useEffect(() => {
    messageBoxRef?.current.scrollIntoView();
  });

  useEffect(() => {
    (async () => {
      try {
        if (router.isReady) {
          const { roomId } = router.query;
          const response = await http.get(`/rooms/${roomId}/messages`, {
            params: {
              limit: 20,
            },
            withCredentials: true,
          });

          setMessages(response.data?.data);
        }
      } catch (e) {
        // TODO: handle failed when fetching data
        console.error(e);
      }

      setIsFetchingMessage(false);
    })();
  }, [isFetchingMessage, router.isReady, router.query]);

  useEffect(() => {
    if (wsInstance) {
      wsInstance.onmessage = (event: MessageEvent) => {
        const response = JSON.parse(event.data);
        setMessages((prevMessages: any) => [response, ...prevMessages]);
      };
    }
  }, [wsInstance]);

  const onSubmitMessage: SubmitHandler<SendMessageInput> = (data) => {
    wsInstance?.send(JSON.stringify(data));
    setIsFetchingMessage(true);
  };

  const getFriendName = (): string => {
    const currentRoom = getCurrentRoom();

    const friend = currentRoom?.participants?.find(
      (participant: any) => participant.id !== currentUser?.id
    );

    return friend?.username;
  };

  return (
    <>
      <div className="flex justify-between items-center sticky top-0 bg-transparent backdrop-blur w-full px-6 py-3 border-b">
        <div className="flex items-center gap-2">
          <Link
            href="/inboxes"
            className="hover:bg-slate-100 hover:rounded-full py-2 p-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth="1.5"
              stroke="currentColor"
              className="w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </Link>

          <p className="font-medium text-lg">{getFriendName()}</p>
        </div>

        <button className="hover:bg-slate-100 hover:rounded-full py-2 p-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z"
            />
          </svg>
        </button>
      </div>

      <div className="flex flex-col min-h-screen flex-col-reverse">
        <form
          onSubmit={handleSubmit(onSubmitMessage)}
          className="flex flex-row items-center gap-4 w-full p-4 bg-white border-t sticky bottom-0"
        >
          <input
            {...register("body")}
            id="body"
            type="text"
            className="w-full rounded-md border-slate-300 bg-white"
            placeholder="Type your message..."
          />
          <button id="attachment" className="text-gray-500">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth="1.5"
              stroke="currentColor"
              className="w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M18.375 12.739l-7.693 7.693a4.5 4.5 0 01-6.364-6.364l10.94-10.94A3 3 0 1119.5 7.372L8.552 18.32m.009-.01l-.01.01m5.699-9.941l-7.81 7.81a1.5 1.5 0 002.112 2.13"
              />
            </svg>
          </button>
          <button
            type="submit"
            className="bg-indigo-500 rounded-full py-2 px-4 hover:bg-indigo-800"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth="1.5"
              stroke="currentColor"
              className="w-6 h-6 text-white"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
              />
            </svg>
          </button>
        </form>

        <ul
          ref={messageBoxRef}
          id="message-list"
          className="flex flex-col-reverse px-4 my-4 gap-3"
        >
          {messages.map(({ id, body, userId, createdAt }: any) => {
            const timestamp = new Date(createdAt).toLocaleTimeString("en-us", {
              timeStyle: "short",
            });

            return currentUser?.id === userId ? (
              <li key={id}>
                <div className="flex justify-end mb-1">
                  <p className="bg-indigo-500 px-4 py-2 rounded-xl text-white">
                    {body}
                  </p>
                </div>
                <div className="flex justify-end">
                  <p className="text-slate-500 text-sm">{timestamp}</p>
                </div>
              </li>
            ) : (
              <li key={id}>
                <div className="flex justify-start mb-1">
                  <p className="bg-slate-200 px-4 py-2 rounded-xl text-black">
                    {body}
                  </p>
                </div>
                <div className="flex justify-start">
                  <p className="text-slate-500 text-sm">{timestamp}</p>
                </div>
              </li>
            );
          })}
        </ul>
      </div>
    </>
  );
};

Inboxes.getLayout = function getLayout(page: ReactElement) {
  return <InboxesLayout>{page}</InboxesLayout>;
};

export default Inboxes;
