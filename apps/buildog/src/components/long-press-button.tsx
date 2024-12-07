import { Button } from "@ui/components/button";
import React, { useState, useRef } from "react";

const colorMap: Record<string, string> = {
  red: "bg-red-500 hover:bg-red-700",
  blue: "bg-blue-500 hover:bg-blue-700",
  green: "bg-green-500 hover:bg-green-700",
  yellow: "bg-yellow-500 hover:bg-yellow-700",
  gray: "bg-gray-500 hover:bg-gray-700",
};

export default function LongPressButton({
  children,
  color,
  onDelete,
}: {
  children: React.ReactNode;
  color: string;
  onDelete: () => void;
}) {
  const [pressing, setPressing] = useState(false);
  const [progress, setProgress] = useState(0);
  const [showError, setShowError] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const pressTimeoutRef = useRef(null);
  const progressIntervalRef = useRef(null);
  const holdDuration = 1000;

  const startPress = () => {
    setPressing(true);
    setProgress(0);
    setShowError(false);

    progressIntervalRef.current = setInterval(() => {
      setProgress((prevProgress) => {
        if (prevProgress >= 100) {
          clearInterval(progressIntervalRef.current);
          return 100;
        }
        return prevProgress + 100 / (holdDuration / 100);
      });
    }, 100);

    pressTimeoutRef.current = setTimeout(() => {
      setSubmitting(true);

      setTimeout(() => {
        setSubmitting(false);
        onDelete();
      }, 2000);

      resetPress();
    }, holdDuration);
  };

  const resetPress = () => {
    clearTimeout(pressTimeoutRef.current);
    clearInterval(progressIntervalRef.current);
    setPressing(false);
    setProgress(0);
  };

  const handleMouseUp = () => {
    if (progress < 50) {
      setShowError(true);
    }
    resetPress();
  };

  const buttonBgClass = colorMap[color] || colorMap["gray"];
  const progressBgClass = `bg-${color}-500`;

  return (
    <div>
      <Button
        className={`relative w-full h-full text-white font-bold rounded-md overflow-hidden ${buttonBgClass} ${
          submitting ? "cursor-not-allowed" : ""
        }`}
        onMouseDown={!submitting ? startPress : undefined}
        onMouseUp={!submitting ? handleMouseUp : undefined}
        onTouchStart={!submitting ? startPress : undefined}
        onTouchEnd={!submitting ? handleMouseUp : undefined}
      >
        {submitting ? (
          <div className="flex justify-center items-center">
            <svg
              className="animate-spin h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              ></circle>
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
              ></path>
            </svg>
          </div>
        ) : (
          <>
            {children}
            <div
              className={`absolute top-0 left-0 h-full ${progressBgClass} opacity-30 transition-all duration-100 ease-linear`}
              style={{ width: `${progress}%` }}
            ></div>
          </>
        )}
      </Button>
      {showError && (
        <p className="text-red-500 mt-2 absolute left-6 bottom-6">
          Hold the button longer to confirm!
        </p>
      )}
    </div>
  );
}
