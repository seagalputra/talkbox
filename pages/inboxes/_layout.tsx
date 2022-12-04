import Link from "next/link";

export default function InboxesLayout({ children }: { children: any }) {
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
            <Link href="/inboxes/9bceabf6ad2a605ea08c2978">
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
            </Link>
          </ul>
        </div>

        <div
          id="sidebar-right"
          className="flex-1 overflow-auto max-h-screen bg-gray-50 scrollbar-hide"
        >
          {children}
        </div>
      </div>
    </main>
  );
}
