import random
import socket
import struct
import os

MAX_IT = 1000
count = [2**i for i in range(3, 10)]

if __name__ == '__main__':

    net_ips = []

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

                ntoa = socket.inet_ntoa(sub_net_ip)
                if ntoa in net_ips:
                    continue
                to_save.append(ntoa)  # append subnet ip
                net_ips.append(ntoa)
                i = 1
                while len(to_save) < pow_c and subnet_count > 0:
                    ip = socket.inet_ntoa(struct.pack(">I", masked_subnet_ip + i))
                    i += 1
                    subnet_count -= 1
                    to_save.append(ip)
            if len(to_save) == len(set(to_save)):
                f.write("\n".join(to_save))
            else:
                to_save = list(set(to_save))


    to_save_prefixes = []

    while len(to_save_prefixes) < MAX_IT:
        gen_pref_ip = net_ips[random.randint(0, len(net_ips) - 1)]
        if gen_pref_ip not in to_save_prefixes:
            to_save_prefixes.append(gen_pref_ip)
        else:
            n_gen = socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))[0:random.randint(1, 9)]
            if n_gen not in to_save_prefixes:
                to_save_prefixes.append(n_gen)

    with open(os.path.join(".", "prefixes", "ip_pref1k.txt"), 'w+') as f:
        f.write("\n".join(to_save_prefixes))

