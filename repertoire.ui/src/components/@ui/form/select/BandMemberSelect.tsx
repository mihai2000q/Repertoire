import {
  Avatar,
  Combobox,
  Group,
  HoverCard,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  Tooltip,
  useCombobox
} from '@mantine/core'
import { BandMember } from '../../../../types/models/Artist.ts'
import { IconUser } from '@tabler/icons-react'
import { useEffect, useState } from 'react'

interface BandMemberSelectProps {
  bandMember: BandMember | null
  setBandMember: (bandMember: BandMember | null) => void
  bandMembers: BandMember[] | undefined
}

function BandMemberSelect({ bandMember, setBandMember, bandMembers }: BandMemberSelectProps) {
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption()
  })

  bandMember = bandMembers ? bandMember : null

  const [value, setValue] = useState<string>(bandMember?.name ?? '')
  const [search, setSearch] = useState(bandMember?.name ?? '')
  useEffect(() => {
    setValue(bandMember?.name ?? '')
    setSearch(bandMember?.name ?? '')
  }, [bandMember])

  const filteredMembers =
    search.trim() !== ''
      ? bandMembers?.filter((member) =>
          member.name.toLowerCase().includes(search.toLowerCase().trim())
        )
      : bandMembers

  const BandMemberHoverCard = () => (
    <HoverCard withArrow={true} openDelay={200} shadow={'md'}>
      <HoverCard.Target>
        <Avatar size={24} color={bandMember.color} src={bandMember.imageUrl} alt={bandMember.name}>
          <IconUser size={15} />
        </Avatar>
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Group gap={'xs'} maw={200} wrap={'nowrap'}>
          <Avatar
            size={'lg'}
            color={bandMember.color}
            src={bandMember.imageUrl}
            alt={bandMember.name}
          >
            <IconUser size={30} />
          </Avatar>
          <Stack gap={0}>
            <Text fw={500} lineClamp={2}>
              {bandMember.name}
            </Text>
            <Text c={'dimmed'} fz={'xs'} lineClamp={1}>
              {bandMember.roles[0].name}
              {bandMember.roles.length > 1 && ' ...'}
            </Text>
          </Stack>
        </Group>
      </HoverCard.Dropdown>
    </HoverCard>
  )

  const BandMemberOption = ({ member }: { member: BandMember }) => (
    <Combobox.Option
      key={member.id}
      value={member.name}
      aria-label={member.name}
      onClick={() => setBandMember(bandMember === member ? null : member)}
    >
      <Group gap={'xs'} wrap={'nowrap'}>
        <Avatar
          size={'sm'}
          variant={'light'}
          color={member.color}
          src={member.imageUrl}
          alt={member.name}
        >
          <IconUser size={14} />
        </Avatar>
        <Text inline fw={500} lineClamp={2}>
          {member.name}
        </Text>
      </Group>
    </Combobox.Option>
  )

  function handleSubmit(valueString: string) {
    setValue(valueString)
    setSearch(valueString)
    combobox.closeDropdown()
  }

  function handleClear() {
    setValue('')
    setSearch('')
    setBandMember(null)
  }

  return (
    <Combobox onOptionSubmit={handleSubmit} store={combobox}>
      <Combobox.Target>
        <Tooltip
          label={'The artist of the song is either not set or not a band'}
          ta={'center'}
          w={190}
          multiline
          position={'top'}
          disabled={bandMembers !== undefined}
        >
          <TextInput
            flex={1}
            label={'Band Member'}
            placeholder={'Choose a member'}
            disabled={bandMembers === undefined}
            leftSection={bandMember ? <BandMemberHoverCard /> : <IconUser size={20} />}
            rightSection={
              bandMember ? <Combobox.ClearButton onClear={handleClear} /> : <Combobox.Chevron />
            }
            value={search}
            onChange={(e) => {
              combobox.openDropdown()
              combobox.updateSelectedOptionIndex()
              setSearch(e.currentTarget.value)
            }}
            onClick={() => combobox.openDropdown()}
            onFocus={() => combobox.openDropdown()}
            onBlur={() => {
              combobox.closeDropdown()
              setSearch(value)
            }}
          />
        </Tooltip>
      </Combobox.Target>

      <Combobox.Dropdown>
        <Combobox.Options>
          <ScrollArea.Autosize mah={200} scrollbarSize={5}>
            {bandMembers?.length === 0 ? (
              <Combobox.Empty>Artist has no members</Combobox.Empty>
            ) : filteredMembers?.length === 0 ? (
              <Combobox.Empty>No members found</Combobox.Empty>
            ) : (
              filteredMembers?.map((bandMember) => (
                <BandMemberOption key={bandMember.id} member={bandMember} />
              ))
            )}
          </ScrollArea.Autosize>
        </Combobox.Options>
      </Combobox.Dropdown>
    </Combobox>
  )
}

export default BandMemberSelect
