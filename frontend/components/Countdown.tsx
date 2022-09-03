import React, { useState, useEffect } from "react";
import {Countdown} from "react-daisyui";

function CountdownThingy({
  hours,
  minutes,
  seconds,
}: {
  hours: number;
  minutes: number;
  seconds: number;
}) {
  const [hoursState, setHours] = useState<number>(hours);
  const [minutesState, setMinutes] = useState<number>(minutes);
  const [secondsState, setSeconds] = useState<number>(seconds);

  // useEffect(() => {
  //   const timer = setTimeout(() => {
  //     setHours((v: number) => (v <= 0 ? hoursState : v - 1));
  //   }, 1000);
  //
  //   return () => {
  //     clearTimeout(timer);
  //   };
  // }, [hoursState]);
  //
  // useEffect(() => {
  //   const timer = setTimeout(() => {
  //     setMinutes((v: number) => (v <= 0 ? minutesState : v - 1));
  //   }, 1000);
  //
  //   return () => {
  //     clearTimeout(timer);
  //   };
  // }, [minutesState]);
  useEffect(() => {
    if (hoursState === 0 && minutesState === 0 && secondsState === 0) {
      return
    }

    if (secondsState === 0 ) {
      if (minutesState !== 0) {
        setMinutes(minutesState - 1)
        setSeconds(59)
      }
    }
    if (minutesState === 0) {
      if (hoursState !== 0) {
        setHours(hoursState - 1)
        setMinutes(59)
      }
    }
    if (hoursState === 0 && minutesState === 0 && secondsState === 0) {
      return
    }
    const timer = setTimeout(() => {
      setSeconds((v: number) => (v <= 0 ? secondsState : v - 1));
    }, 1000);

    return () => {
      clearTimeout(timer);
    };
  }, [secondsState]);

  return (
    <div className="grid grid-flow-col gap-5 text-center auto-cols-max">
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={hoursState} />
        hours
      </div>
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={minutesState} />
        min
      </div>
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={secondsState} />
        sec
      </div>
    </div>
  );
}

export default CountdownThingy;
