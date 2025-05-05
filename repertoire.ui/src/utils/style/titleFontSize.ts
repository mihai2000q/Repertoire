const breakpoints = {
  xl: 50,
  lg: 36,
  md: 26,
  sm: 20,
  xs: 16,
  xxs: 12
}

export default function titleFontSize(str: string | undefined | null): string {
  if (str === undefined || str === null || str.trim() === '') return '0'

  const size =
    str.length > breakpoints.xl
      ? { viewport: 2.5, pixels: 28 }
      : str.length > breakpoints.lg
        ? { viewport: 2.55, pixels: 30 }
        : str.length > breakpoints.md
          ? { viewport: 3.1, pixels: 33 }
          : str.length > breakpoints.sm
            ? { viewport: 3.55, pixels: 33 }
            : str.length > breakpoints.xs
              ? { viewport: 3.9, pixels: 35 }
              : str.length > breakpoints.xxs
                ? { viewport: 4.2, pixels: 38 }
                : { viewport: 4.5, pixels: 48 }
  return `max(${size.viewport}vw, ${size.pixels}px)`
}
