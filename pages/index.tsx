import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Link from "next/link";
import http from "../lib/http";
import { AxiosError } from "axios";
import { useRouter } from "next/router";

type UserLoginInput = {
  username: string;
  password: string;
};

type UserLoginErrorOutput = {
  status?: string;
  message?: string;
};

export default function Home() {
  const [errorResponse, setErrorResponse] = useState<UserLoginErrorOutput>();
  const { register, handleSubmit } = useForm<UserLoginInput>();
  const router = useRouter();

  const onSubmitLogin: SubmitHandler<UserLoginInput> = async (data) => {
    try {
      await http.post("/auth/login", data, {
        withCredentials: true,
      });

      router.push("/inboxes");
    } catch (err) {
      if (err instanceof AxiosError) {
        const errResponse = err.response?.data;
        setErrorResponse(errResponse);
      }
    }
  };

  return (
    <main className="container mx-auto">
      <div className="flex flex-col justify-center items-center min-h-screen">
        <form
          onSubmit={handleSubmit(onSubmitLogin)}
          className="w-1/2 lg:w-1/4 flex flex-col gap-3"
        >
          <div className="flex bg-rose-400 text-white">
            <p>{errorResponse?.message}</p>
          </div>
          <div className="flex flex-col gap-2">
            <label htmlFor="username" className="text-slate-500">
              Username
            </label>
            <div className="flex flex-col gap-1">
              <input
                {...register("username")}
                id="username"
                type="text"
                className="w-full rounded-md border-slate-300 bg-white"
                placeholder="Enter your email address..."
              />
            </div>
          </div>

          <div className="flex flex-col gap-2">
            <label htmlFor="password" className="text-slate-500">
              Password
            </label>
            <div className="flex flex-col gap-1">
              <input
                {...register("password")}
                id="password"
                type="password"
                className="w-full rounded-md border-slate-300 bg-white"
                placeholder="Enter your password..."
              />
            </div>
          </div>

          <button
            type="submit"
            className="w-full text-center bg-indigo-500 text-white font-medium rounded-md py-2 hover:bg-indigo-700 hover:cursor-pointer"
          >
            Submit
          </button>
        </form>

        <div className="w-1/2 lg:w-1/3 mt-12">
          <p className="text-slate-500 text-xs text-center">
            By clicking the submit button above, you acknowledge that you have
            read and understood our{" "}
            <Link href="#" className="text indigo-500 underline">
              Terms & Conditions
            </Link>{" "}
            and{" "}
            <Link href="#" className="text indigo-500 underline">
              Privacy and Policy
            </Link>
            .
          </p>
        </div>
      </div>
    </main>
  );
}
