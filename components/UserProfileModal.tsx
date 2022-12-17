import { AxiosError } from "axios";
import { useState, useEffect, ChangeEvent } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { useCookies } from "react-cookie";
import { useRouter } from "next/router";
import Image from "next/image";
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
  const [errorResponse, setErrorResponse] = useState<UpdateProfileOutput>();
  const [userProfile, setUserProfile] = useState<UserProfile>({});
  const { register, handleSubmit, setValue } = useForm<UpdateUserInput>();
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);
  const [, , removeCookie] = useCookies(["talkbox"]);
  const router = useRouter();

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

  const onSubmitAvatar = async (event: ChangeEvent<HTMLInputElement>) => {
    event.preventDefault();
    const avatarFile: File | null | undefined = event.target.files?.item(0);

    const formData = new FormData();
    if (avatarFile) {
      formData.append("avatar", avatarFile);
    }

    try {
      const response = await http.post("/users/avatar", formData, {
        withCredentials: true,
      });
      const data = response.data?.data;

      const imageUrl = data?.imageUrl;
      setUserProfile((prevData) => ({
        ...prevData,
        avatar: imageUrl,
      }));
    } catch (err) {
      if (err instanceof AxiosError) {
        const errResponse = err.response?.data;
        setErrorResponse(errResponse);
      }
    }
  };

  const visiblePassword = (event: any) => {
    event.preventDefault();
    setIsPasswordVisible((prevValue) => !prevValue);
  };

  const onUserLogout = (event: any) => {
    event.preventDefault();
    removeCookie("talkbox");
    router.push("/");
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

        {errorResponse ? (
          <div className="flex bg-rose-400 text-white mt-4 p-4 rounded">
            <p>{errorResponse?.message}</p>
          </div>
        ) : null}

        <div className="mt-4 relative flex justify-center">
          <Image
            src={userProfile?.avatar || ""}
            className="rounded-full w-24"
            alt="User avatar"
            width={200}
            height={200}
          />

          <div className="rounded-full w-24 absolute top-0 w-full h-full flex items-center justify-center text-transparent hover:bg-black/[0.4] hover:text-white transition duration-200">
            <input
              type="file"
              onChange={onSubmitAvatar}
              className="opacity-0 absolute"
            />
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
            <div className="flex border border-slate-300 rounded-md shadow-sm mt-1">
              <input
                {...register("password")}
                type={isPasswordVisible ? "text" : "password"}
                name="password"
                id="password"
                className="flex-1 border-transparent bg-white rounded-md"
              />
              <button className="mx-3 text-gray-500" onClick={visiblePassword}>
                {isPasswordVisible ? (
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
                      d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
                    />
                  </svg>
                ) : (
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
                      d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                  </svg>
                )}
              </button>
            </div>
          </div>
          <button
            type="submit"
            className="col-span-2 text-center bg-indigo-500 text-white font-medium rounded-md py-2 hover:bg-indigo-700 hover:cursor-pointer"
          >
            Submit
          </button>
        </form>

        <button
          type="button"
          onClick={onUserLogout}
          className="w-full text-center text-red-500 font-medium rounded-md py-2 mt-2 hover:cursor-pointer"
        >
          Logout
        </button>
      </div>
    </div>
  );
};

export default UserProfileModal;
