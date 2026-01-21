import "@mantine/core";

import type { EmotionStyles, EmotionSx } from "@mantine/emotion";

declare module "@mantine/core" {
  // noinspection JSUnusedGlobalSymbols
  export interface BoxProps {
    sx?: EmotionSx;
    styles?: EmotionStyles;
  }
}
