import arcpy

from datetime import date, timedelta

# ___Config___

# Shapefile path
DATAPATH = "../Data/DatePoints/"
arcpy.env.workspace = DATAPATH

# Shapefile name
POINT_DATA = "points.shp"

# Name of date attribute field
DATE_FIELD_NAME = "Date"

# Max distance in meters
D_MAX = 1800
# Max temporal difference in days
T_MAX = 16


def main():
    Ns, N2s, Nt, N2t, X, n = calcSpaceTimeCluster()

    print("\nCounts:")
    print("Ns: " + str(Ns))
    print("N2s: " + str(N2s))
    print("Nt: " + str(Nt))
    print("N2t: " + str(N2t))
    print("X: " + str(X))
    print("n: " + str(n))

    N, E, V = calc_statistics(Ns, N2s, Nt, N2t, X, n)

    print("\nStatistics:")
    print("N: " + str(N))
    print("E: " + str(E))
    print("V: " + str(V))


def dateFieldValid():
    fields = arcpy.ListFields(POINT_DATA)
    for field in fields:
        if field.name == DATE_FIELD_NAME and field.type == 'Date':
            print("Valid date-field.")
            return True

    print("Date field does not exist. Aborting..")
    return False

def calcSpaceTimeCluster():
    print("====================================================")
    if not dateFieldValid():
        return

    print("Init space time cluster calculation")
    fields = ["FID", DATE_FIELD_NAME, "SHAPE@"]

    # Distance matches (first and second order)
    Ns = 0
    N2s = 0
    # Time matches (first and second order)
    Nt = 0
    N2t = 0
    # Both matching
    X = 0

    # Num points
    n = 0


    with arcpy.da.SearchCursor(POINT_DATA, fields) as i_cursor:
        for i_point in i_cursor:
            n += 1
            if n % 10 == 0:
                print(str(n) + " features complete")
            with arcpy.da.SearchCursor(POINT_DATA, fields) as j_cursor:
                for j_point in j_cursor:
                    # Do not count i==j matches
                    if i_point[0] == j_point[0]:
                        continue

                    d_match = False
                    t_match = False

                    # Check if below distance threshold
                    if distance_diff(i_point[2], j_point[2]) <= D_MAX:
                        d_match = True
                        Ns += 1

                    # Check if below temporal threshold
                    if time_diff(i_point[1], j_point[1]) <= T_MAX:
                        t_match = True
                        Nt += 1

                    if d_match and t_match:
                        X += 1

                    # Second order counting match j on k
                    with arcpy.da.SearchCursor(POINT_DATA, fields) as k_cursor:
                        for k_point in k_cursor:
                            if i_point[0] == k_point[0] or j_point[0] == k_point[0]:
                                continue

                            if d_match:
                                if distance_diff(j_point[2], k_point[2]) <= D_MAX:
                                    N2s += 1

                            if t_match:
                                if time_diff(j_point[1], k_point[1]) <= T_MAX:
                                    N2t += 1

    # Normalize for double counting and return
    return normalize_double_count(Ns, N2s, Nt, N2t, X, n)


def distance_diff(point1_geometry, point2_geometry):
    dist = point1_geometry.distanceTo(point2_geometry)
    return dist


def time_diff(date1, date2):
    delta = date2 - date1
    return abs(delta.days)


def normalize_double_count(Ns, N2s, Nt, N2t, X, n):
    return Ns/2, N2s/2, Nt/2, N2t/2, X/2, n


def calc_statistics(Ns, N2s, Nt, N2t, X, n):
    # Number of pairs
    N = n * (n - 1) / 2
    # Expected value
    E = Nt * Ns / N
    # Variance
    V = Ns*Nt/N + 4*N2s*N2t/(n*(n - 1)*(n - 2)) + \
        4*(Ns*(Ns - 1) - N2s)*(Nt*(Nt - 1) - N2t)/(n*(n - 1)*(n - 2)*(n - 3)) - \
        (Ns*Nt/N)*(Ns*Nt/N)

    return N, E, V

main()



