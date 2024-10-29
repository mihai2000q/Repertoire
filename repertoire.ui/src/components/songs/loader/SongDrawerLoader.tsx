import {Skeleton, Stack} from "@mantine/core";

function SongDrawerLoader() {
  return (
    <Stack gap={'xs'} data-testid={'song-drawer-loader'}>
      <Skeleton radius={0} w={'100%'} h={400} />

      <Stack gap={'xs'} px={'md'}>
        <Skeleton w={300} h={25} />
        <Stack gap={2}>
          <Skeleton w={400} h={8} />
          <Skeleton w={400} h={8} />
          <Skeleton w={50} h={8} />
        </Stack>
      </Stack>
    </Stack>
  );
}

export default SongDrawerLoader;
