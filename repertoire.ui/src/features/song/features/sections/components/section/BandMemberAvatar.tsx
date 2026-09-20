import { Avatar, AvatarProps, Group, HoverCard, Stack, Text } from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import { BandMember } from '../../../../../../types/models/Artist.ts'

interface BandMemberAvatarProps extends AvatarProps {
  bandMember: BandMember
  iconSize?: string | number
}

function BandMemberAvatar({ bandMember, iconSize = 15, ...props }: BandMemberAvatarProps) {
  return (
    <HoverCard openDelay={200} position="top">
      <HoverCard.Target>
        <Avatar
          size={25}
          {...props}
          color={bandMember.color}
          src={bandMember.imageUrl}
          alt={bandMember.imageUrl && bandMember.name}
        >
          <IconUser aria-label={`default-icon-${bandMember.name}`} size={iconSize} />
        </Avatar>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} wrap={'nowrap'}>
          <Avatar
            size={60}
            color={bandMember.color}
            src={bandMember.imageUrl}
            alt={bandMember.imageUrl && bandMember.name}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          >
            <IconUser size={30} />
          </Avatar>
          <Stack gap={0}>
            <Text fw={500} lineClamp={2}>
              {bandMember.name}
            </Text>
            {bandMember.roles.slice(0, 2).map((role, index) => (
              <Text key={role.id} c={'dimmed'} fz={'xs'} lineClamp={1} lh={1.05}>
                {role.name}
                {index === 1 && bandMember.roles.length > 2 && ' ...'}
              </Text>
            ))}
          </Stack>
        </Group>
      </HoverCard.Dropdown>
    </HoverCard>
  )
}

export default BandMemberAvatar
