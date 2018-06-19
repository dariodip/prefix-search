import random
import socket
import struct
import os

count = [2**i for i in range(3, 10)]

if __name__ == '__main__':

    for pow_c in count:
        to_save = []
        with open(os.path.join('.', 'dataset', 'ip{}.txt'.format(pow_c)), "w+") as f:
            while len(to_save) < pow_c:
                s_gen_ip = socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))
                gen_ip = socket.inet_aton(s_gen_ip)
                subnet_count = random.randint(0, (pow_c - 1) % 32)
                subnet_ip_unpacked = struct.unpack(">I", gen_ip)[0]
                mask = (2**32)-1 & ~((2 ** subnet_count) - 1)
                masked_subnet_ip = subnet_ip_unpacked & mask
                sub_net_ip = struct.pack(">I", masked_subnet_ip)

                to_save.append(socket.inet_ntoa(sub_net_ip))  # append subnet ip

                i = 1
                while len(to_save) < pow_c and subnet_count > 0:
                    ip = socket.inet_ntoa(struct.pack(">I", masked_subnet_ip + i))
                    i += 1
                    subnet_count -= 1
                    to_save.append(ip)

            f.write("\n".join(to_save))
