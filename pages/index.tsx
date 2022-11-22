import Link from "next/link";

export default function Home() {
  return (
    <main className="container mx-auto">
      <div className="flex flex-col justify-center items-center min-h-screen">
        <div className="w-1/4 flex flex-col gap-3">
          <div className="flex flex-col gap-2">
            <label htmlFor="email" className="text-slate-500">
              Email
            </label>
            <div className="flex flex-col gap-1">
              <input
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
                type="password"
                className="w-full rounded-md border-slate-300 bg-white"
                placeholder="Enter your password..."
              />
            </div>
          </div>

          <Link
            href="/inboxes"
            className="w-full text-center bg-indigo-500 text-white font-medium rounded-md py-2 hover:bg-indigo-700 hover:cursor-pointer"
          >
            Submit
          </Link>
        </div>

        <div className="w-1/3 mt-12">
          <p className="text-slate-500 text-xs text-center">
            By clicking the submit button above, you acknowledge that you have
            read and understood our
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
