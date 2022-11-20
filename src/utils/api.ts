import Joi from "joi";

enum ResponseStatus {
  success = "success",
  error = "error",
}

const headers = {
  "Content-Type": "application/json",
};

const defaultErrorMessage = "Your request can't be processed, please try again";

function successResponse(body: object, meta: any = {}): Response {
  const data = JSON.stringify({
    status: ResponseStatus.success,
    meta,
    data: body,
  });
  return new Response(data, {
    status: 200,
    headers,
  });
}

function errorResponse(
  error?: any | null,
  statusCode: number = 422,
  message: string = defaultErrorMessage
): Response {
  const response = {
    status: ResponseStatus.error,
    message,
  };
  if (error instanceof Joi.ValidationError) {
    const errorDetails = error.details.map((detail) => ({
      field: detail?.path[0],
      message: detail?.message,
    }));

    return new Response(
      JSON.stringify({
        ...response,
        message: errorDetails,
      }),
      { status: 400, headers }
    );
  }

  const data = JSON.stringify(response);

  return new Response(data, {
    status: statusCode,
    headers,
  });
}

export { successResponse, errorResponse, ResponseStatus };
