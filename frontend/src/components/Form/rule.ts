import { mergeRight } from "ramda";

export const requiredRule = {
  required: "Required",
};

export const userNameRule = mergeRight(requiredRule, {
  minLength: {
    value: 3,
    message: "Username must be at least 3 characters",
  },
  maxLength: {
    value: 20,
    message: "Username must be at most 20 characters",
  },
  pattern: {
    value: /^[a-zA-Z0-9_]+$/,
    message: "Username can only contain letters, numbers and underscores",
  },
});

export const passwordRule = mergeRight(requiredRule, {
  minLength: {
    value: 6,
    message: "The password must be at least 6 characters long",
  },
  maxLength: {
    value: 20,
    message: "The password must be at most 20 characters long",
  },
  pattern: {
    value: /^[a-zA-Z0-9_]+$/,
    message: "Password can only contain letters, numbers and underscores",
  },
});

export const amountRule = mergeRight(requiredRule, {
  pattern: {
    value: /^[0-9]+(\.[0-9]{1,6})?$/,
    message: "Please enter a legal amount, with up to 6 decimal places.",
  },
  min: {
    value: 0.000001,
    message: "Amount must be greater than 0",
  },
});
