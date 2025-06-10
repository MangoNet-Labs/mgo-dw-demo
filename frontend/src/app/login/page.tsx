"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useLoginEffect } from "@/hooks/useLoginEffect";
import { useForm } from "react-hook-form";
import Image from "next/image";
import { loginEffect, LoginProps } from "@/effector/effector";
import { useUnit } from "effector-react";
import Link from "next/link";
import useSignedGuard from "@/hooks/useSignedGuard";
import { passwordRule, userNameRule } from "@/components/Form/rule";
import { path } from "ramda";
export default function Home() {
  const { handel } = useLoginEffect();
  const loginPending = useUnit(loginEffect.pending);
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginProps>();
  
  useSignedGuard();
  return (
    <form onSubmit={handleSubmit(handel)} className="space-y-4 p-4 h-full">
      <div className="h-full w-[93%] m-auto flex flex-col items-baseline justify-around">
        <div className="w-full pt-24">
          <div className="w-[80%] m-auto relative h-11">
            <Image
              src={"/images/mgonetwork.png"}
              fill
              className="object-contain"
              alt={""}
            />
          </div>
          <p className="text-xl text-center mt-2 text-[var(--t-main)]">
            Deposit and Withdrawal DEMO
          </p>
        </div>

        <div className="w-full">
          <p className="text-xl text-center mb-7 text-white">Log in</p>

          <div className="flex flex-col gap-4 m-auto">
            <Input
              {...register("username")}
              label="username"
              placeholder="Please enter"
              {...register("username", userNameRule)}
              error={path(["username", "message"], errors)}
            />

            <Input
              {...register("password")}
              label="Password"
              type="password"
              placeholder="Please enter"
              {...register("password", passwordRule)}
              error={path(["password", "message"], errors)}
            />
          </div>
        </div>

        <div className="w-full">
          <Button
            isLoading={loginPending}
            type="submit"
            className="bg-[var(--t-main)] w-full"
          >
            Log in
          </Button>
          <h4 className="text-white text-xs text-center mt-4 flex items-center justify-center gap-0.5">
            <span>No account?</span>
            <Link className="text-[var(--t-main)]" href={"register"}>
              Go to register
            </Link>
          </h4>
        </div>
      </div>
    </form>
  );
}
