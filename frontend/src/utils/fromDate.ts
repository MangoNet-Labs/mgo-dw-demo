import moment from "moment";

export const formatDateTime = (date: string) => {
  if (!date) return "--";
  if (date.length < 10) return "--";
  return moment.unix(parseInt(date)).format("YYYY-MM-DD");
};

export const formatTime = (seconds: number) => {
  console.log(seconds);
  // const days = Math.floor(seconds / (24 * 60 * 60));
  // const hours = Math.floor((seconds % (24 * 60 * 60)) / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = (seconds % 60).toFixed(0);
  return `${mins ? mins : "0"} M : ${secs} S`;
};

export const calculateDays = (startTime: string, endTime: string) => {
  const date1 = moment.unix(parseInt(startTime));
  const date2 = moment.unix(parseInt(endTime));

  const diffDays = date2.diff(date1, "days");
 
  return diffDays;
};
