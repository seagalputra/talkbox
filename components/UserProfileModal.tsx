type UpdateUserInput = {
  firstName?: string;
  lastName?: string;
  username?: string;
  password?: string;
  avatar?: string;
};

type UpdateUserOutput = {
  status?: string;
  message?: string;
};

const UserProfileModal = ({ openUserProfileModal }: any) => {
  return (
    <div className="fixed z-10 top-0 left-0 w-full h-full overflow-auto bg-black/[0.4] flex items-center justify-center">
      <div className="bg-white mx-auto border w-1/3 rounded p-6">
        <div className="flex justify-between">
          <h1 className="font-bold text-3xl text-slate-800">User Profile</h1>
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

        <div className="flex justify-center mt-4">
          <img
            className="rounded-full w-24 border-2 border-slate-100"
            src="https://i.picsum.photos/id/524/200/200.jpg?hmac=t6LNfKKZ41wUVh8ktcFHag3CGQDzovGpZquMO5cbH-o"
            alt="User avatar"
          />
        </div>
        <form className="mt-4 grid grid-cols-2 gap-6">
          <div className="col-span-1">
            <label
              htmlFor="firstName"
              className="block font-medium text-gray-700 text-sm"
            >
              First Name
            </label>
            <input
              type="text"
              name="firstName"
              id="firstName"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
            />
          </div>
          <div className="col-span-1">
            <label
              htmlFor="lastName"
              className="block font-medium text-gray-700 text-sm"
            >
              Last Name
            </label>
            <input
              type="text"
              name="lastName"
              id="lastName"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
            />
          </div>
          <div className="col-span-2">
            <label
              htmlFor="email"
              className="block font-medium text-gray-700 text-sm"
            >
              Email
            </label>
            <input
              type="text"
              name="email"
              id="email"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
            />
          </div>
          <div className="col-span-2">
            <label
              htmlFor="username"
              className="block font-medium text-gray-700 text-sm"
            >
              Username
            </label>
            <input
              type="text"
              name="username"
              id="username"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
            />
          </div>
          <div className="col-span-2">
            <label
              htmlFor="password"
              className="block font-medium text-gray-700 text-sm"
            >
              Password
            </label>
            <input
              type="text"
              name="password"
              id="password"
              className="w-full rounded-md border-slate-300 bg-white mt-1 shadow-sm"
            />
          </div>
          <button className="col-span-2 text-center bg-indigo-500 text-white font-medium rounded-md py-2 hover:bg-indigo-700 hover:cursor-pointer">
            Submit
          </button>
        </form>
      </div>
    </div>
  );
};

export default UserProfileModal;
