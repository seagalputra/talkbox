const NewMessageModal = ({ openNewMessageModal }: any) => {
  return (
    <div className="fixed z-10 top-0 left-0 w-full h-full overflow-auto bg-black/[0.4] flex items-center justify-center">
      <div className="bg-white mx-auto border w-1/3 rounded p-6">
        <div className="flex justify-between">
          <h1 className="font-bold text-3xl text-slate-800">New Message</h1>
          <button onClick={openNewMessageModal}>
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
          </button>
        </div>

        <form className="mt-4 grid grid-cols-2 gap-6">
          <div className="col-span-2">
            <input
              type="text"
              name="email"
              id="email"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
              placeholder="Enter friend email address"
            />
          </div>
        </form>
      </div>
    </div>
  );
};

export default NewMessageModal;
