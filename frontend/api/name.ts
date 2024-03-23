import { useUserStore } from "@/stores/user"

const handlePacket = (data: string) => {
  const user = useUserStore()
  if (user.data) {
    user.data.name = data
  }
}

export default {handlePacket}