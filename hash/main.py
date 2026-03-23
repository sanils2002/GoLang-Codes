# 
def convert_ip_range(row_key):
    if '/' in row_key:
        parts = row_key.split("/")
        print(parts)

        network_address = parts[0] + ":" + parts[1]
        return network_address
    return row_key

# Example usage:
row_key = "192.168.1.0"
print(convert_ip_range(row_key))