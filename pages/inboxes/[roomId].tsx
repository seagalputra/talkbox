import Link from "next/link";
import { useRouter } from "next/router";
import { SubmitHandler, useForm } from "react-hook-form";
import { useEffect, useState } from "react";

type SendMessageInput = {
  body: string;
  attachment?: string;
};

export default function Inboxes() {
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
    <main className="container mx-auto">
      <div className="flex w-full bg-white min-h-screen divide-x border-x">
        <div
          id="sidebar-left"
          className="w-1/3 flex flex-col max-h-screen overflow-auto"
        >
          <div className="flex items-center w-full justify-between px-6 sticky top-0 bg-transparent backdrop-blur">
            <h1 className="font-bold text-3xl text-slate-800 py-4">Inbox</h1>
            <div className="flex gap-2 items-center">
              <button className="hover:bg-slate-100 hover:rounded-full py-2 p-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="w-6 h-6 text-indigo-500"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                  />
                </svg>
              </button>
              <button className="hover:bg-slate-100 hover:rounded-full py-2 p-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  className="w-6 h-6 text-indigo-500"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
              </button>
            </div>
          </div>

          <ul id="inbox-list" className="flex flex-col mt-2 divide-y">
            <li className="flex flex-row gap-4 p-4 hover:bg-slate-100 hover:cursor-pointer">
              <img
                className="avatar rounded-full w-16"
                src="https://i.picsum.photos/id/524/200/200.jpg?hmac=t6LNfKKZ41wUVh8ktcFHag3CGQDzovGpZquMO5cbH-o"
                alt="User avatar"
              />
              <div className="flex flex-row justify-between w-full">
                <div className="flex flex-col gap-2">
                  <p className="font-bold font-sans text-md text-slate-800 mt-1">
                    John Doe
                  </p>
                  <p className="font-sans text-slate-400 text-sm">
                    Hello World!
                  </p>
                </div>

                <div className="flex flex-col gap-2">
                  <p className="text-sm text-slate-400 mt-1">11.00</p>
                  <p className="text-sm border rounded-full text-center bg-red-500 border-red-500 text-white">
                    1
                  </p>
                </div>
              </div>
            </li>
          </ul>
        </div>

        <div
          id="sidebar-right"
          className="flex-1 overflow-auto max-h-screen bg-gray-50 scrollbar-hide"
        >
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
        </div>
      </div>
    </main>
  );
}
