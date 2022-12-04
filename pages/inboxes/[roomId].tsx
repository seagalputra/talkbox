import Link from "next/link";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import type { NextPageWithLayout } from "../_app";
import type { ReactElement } from "react";
import InboxesLayout from "./_layout";

type SendMessageInput = {
  body: string;
  attachment?: string;
};

const Inboxes: NextPageWithLayout = () => {
  const router = useRouter();
  const [wsInstance, setWsInstance] = useState<any>(null);
  const { register, handleSubmit, formState, reset } =
    useForm<SendMessageInput>();
  const [messagesEnd, setMessagesEnd] = useState<any>(null);

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
    messagesEnd?.scrollIntoView({ behavior: "smooth" });
  }, [messagesEnd]);

  const onSubmitMessage: SubmitHandler<SendMessageInput> = (data) => {
    wsInstance?.send(JSON.stringify(data));
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

          <p className="font-medium text-lg">John Doe</p>
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
          id="message-list"
          className="flex flex-col px-4 mb-4"
          ref={(element) => {
            setMessagesEnd(element);
          }}
        >
          <li>
            <div className="flex justify-start mb-1">
              <p className="bg-slate-200 px-4 py-2 rounded-xl text-black">
                Hello World!
              </p>
            </div>
            <div className="flex justify-start">
              <p className="text-slate-500 text-sm">10.50</p>
            </div>
          </li>
        </ul>
      </div>
    </>
  );
};

Inboxes.getLayout = function getLayout(page: ReactElement) {
  return <InboxesLayout>{page}</InboxesLayout>;
};

export default Inboxes;
