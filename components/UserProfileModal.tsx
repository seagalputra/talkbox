const UserProfileModal = ({ openUserProfileModal }: any) => {
  return (
    <div className="fixed z-10 top-0 left-0 w-full h-full overflow-auto bg-black/[0.4]">
      <div className="bg-white my-[15%] mx-auto border w-1/3 rounded p-4">
        <div className="flex justify-between">
          <p className="text-slate-700 font-bold">User Profile</p>
          <button onClick={openUserProfileModal}>
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
      </div>
    </div>
  );
};

export default UserProfileModal;
