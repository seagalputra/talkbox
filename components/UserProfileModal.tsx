import { AxiosError } from "axios";
import { useState, useEffect } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import http from "../lib/http";

type UpdateUserInput = {
  firstName?: string;
  lastName?: string;
  username?: string;
  email?: string;
  password?: string;
  avatar?: string;
};

type UserProfile = UpdateUserInput & {
  createdAt?: string;
  updatedAt?: string;
};

type UpdateProfileOutput = {
  status?: string;
  message?: string;
};

const UserProfileModal = ({ openUserProfileModal }: any) => {
  const [errorResponse, setErrorResponse] = useState<UpdateProfileOutput>({});
  const [userProfile, setUserProfile] = useState<UserProfile>({});
  const { register, handleSubmit, setValue } = useForm<UpdateUserInput>();

  useEffect(() => {
    (async () => {
      try {
        const response = await http.get("/users/profile", {
          withCredentials: true,
        });

        const data: UserProfile = response?.data.data;

        setValue("firstName", data?.firstName);
        setValue("lastName", data?.lastName);
        setValue("email", data?.email);
        setValue("username", data?.username);

        setUserProfile(data);
      } catch (e) {
        // TODO: handle error when failed to fetch user profile
        console.error(e);
      }
    })();
  }, [setValue]);

  const onSubmitUserProfile: SubmitHandler<UpdateUserInput> = async (data) => {
    try {
      await http.patch("/users", data, {
        withCredentials: true,
      });

      openUserProfileModal();
    } catch (err) {
      if (err instanceof AxiosError) {
        const errResponse = err.response?.data;
        setErrorResponse(errResponse);
      }
    }
  };

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

        <div className="mt-4 relative flex justify-center">
          <img
            className="rounded-full w-24"
            src="https://i.picsum.photos/id/524/200/200.jpg?hmac=t6LNfKKZ41wUVh8ktcFHag3CGQDzovGpZquMO5cbH-o"
            alt="User avatar"
          />

          <div className="rounded-full w-24 absolute top-0 w-full h-full flex items-center justify-center text-transparent hover:bg-black/[0.4] hover:text-white transition duration-200">
            <input type="file" className="opacity-0 absolute" />
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
                d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z"
              />
            </svg>
          </div>
        </div>
        <form
          onSubmit={handleSubmit(onSubmitUserProfile)}
          className="mt-4 grid grid-cols-2 gap-6"
        >
          <div className="col-span-1">
            <label
              htmlFor="firstName"
              className="block font-medium text-gray-700 text-sm"
            >
              First Name
            </label>
            <input
              {...register("firstName")}
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
              {...register("lastName")}
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
              {...register("email")}
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
              {...register("username")}
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
              {...register("password")}
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
