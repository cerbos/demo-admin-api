export {};

declare global {
  interface Window {
    MonacoEnvironment: any; // 👈️ turn off type checking
  }
}
