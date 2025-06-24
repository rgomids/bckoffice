"use client";
import { Transition } from "@headlessui/react";
import { motion, AnimatePresence } from "framer-motion";

interface ToastProps {
  message: string | null;
}

export default function Toast({ message }: ToastProps) {
  return (
    <AnimatePresence>
      {message && (
        <Transition
          show={true}
          as={motion.div}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 20 }}
          className="fixed bottom-4 right-4 bg-card text-background px-4 py-2 rounded"
        >
          {message}
        </Transition>
      )}
    </AnimatePresence>
  );
}
