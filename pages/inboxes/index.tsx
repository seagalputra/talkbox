import type { NextPageWithLayout } from "../_app";
import type { ReactElement } from "react";
import InboxesLayout from "./_layout";

const Inboxes: NextPageWithLayout = () => {
  return (
    <div className="flex justify-center items-center h-full">
      <p>Select a chat to start messaging</p>
    </div>
  );
};

Inboxes.getLayout = function getLayout(page: ReactElement) {
  return <InboxesLayout>{page}</InboxesLayout>;
};

export default Inboxes;
