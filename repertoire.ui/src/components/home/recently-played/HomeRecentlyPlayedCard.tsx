import {
  alpha,
  Avatar,
  Center,
  Grid,
  Group,
  GroupProps,
  Stack,
  Text,
  Tooltip,
  useMatches
} from '@mantine/core'
import ProgressBar from '../../@ui/bar/ProgressBar.tsx'
import dayjs from 'dayjs'
import PastDate from '../../@ui/text/PastDate.tsx'
import { useHover, useMergedRef } from '@mantine/hooks'
import { MouseEvent, ReactNode, RefObject } from 'react'

interface AdditionalTextProps {
  content: string
  onClick: (e: MouseEvent) => void
}

interface HomeRecentlyPlayedCardProps extends GroupProps {
  ref: RefObject<HTMLDivElement>
  imageUrl: string | null | undefined
  title: string
  progress: number
  lastPlayed: string | null | undefined
  openedMenu: boolean
  defaultIcon: ReactNode
  onClick: () => void
  isArtist?: boolean
  additionalText?: AdditionalTextProps
}

function HomeRecentlyPlayedCard({
  ref: refProp,
  imageUrl,
  title,
  progress,
  lastPlayed,
  openedMenu,
  defaultIcon,
  isArtist,
  onClick,
  additionalText,
  ...others
}: HomeRecentlyPlayedCardProps) {
  const { ref: hoverRef, hovered } = useHover()
  const ref = useMergedRef(refProp, hoverRef)

  const isSelected = hovered || openedMenu === true

  // Not Recommended usage
  const groupGap = useMatches({
    base: 'md',
    lg: 'xs',
    xl: 'md'
  })

  return (
    <Group
      ref={ref}
      wrap={'nowrap'}
      sx={(theme) => ({
        transition: '0.3s',
        border: '1px solid transparent',
        ...(isSelected && {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.primary[0], 0.15)
        })
      })}
      pl={'lg'}
      pr={'xxs'}
      py={'xs'}
      gap={groupGap}
      onClick={onClick}
      {...others}
    >
      <Avatar
        radius={isArtist === true ? '50%' : 'md'}
        src={imageUrl}
        alt={imageUrl && title}
        bg={isArtist === true ? 'gray.0' : 'gray.5'}
        sx={(theme) => ({
          aspectRatio: 1,
          boxShadow: theme.shadows.sm
        })}
      >
        <Center c={isArtist === true ? 'gray.7' : 'white'}>{defaultIcon}</Center>
      </Avatar>

      <Grid flex={1} columns={12} align={'center'}>
        <Grid.Col span={{ base: 5, md: 8, xxl: 5 }}>
          <Stack gap={0} style={{ overflow: 'hidden' }}>
            <Text fw={additionalText ? 600 : 500} lineClamp={1}>
              {title}
            </Text>
            {additionalText && (
              <Group>
                <Text
                  fz={'sm'}
                  c={'dimmed'}
                  lineClamp={1}
                  sx={{ '&:hover': { textDecoration: 'underline' } }}
                  style={{ cursor: 'pointer' }}
                  onClick={additionalText.onClick}
                >
                  {additionalText.content}
                </Text>
              </Group>
            )}
          </Stack>
        </Grid.Col>
        <Grid.Col span={4} display={{ base: 'block', md: 'none', xxl: 'block' }}>
          <ProgressBar progress={progress} mx={'xs'} />
        </Grid.Col>
        <Grid.Col span={{ base: 3, md: 4, xxl: 3 }} px={'md'}>
          <Tooltip
            label={`Last time played on ${dayjs(lastPlayed).format('D MMMM YYYY [at] hh:mm A')}`}
            openDelay={400}
            disabled={!lastPlayed}
          >
            <PastDate
              dateValue={lastPlayed}
              ta={'center'}
              fz={'sm'}
              fw={500}
              c={'dimmed'}
              truncate={'end'}
            />
          </Tooltip>
        </Grid.Col>
      </Grid>
    </Group>
  )
}

export default HomeRecentlyPlayedCard
