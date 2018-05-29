import arcpy
import math
import time
import json
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


# Export data for other implementaions
EXPORT = True
EXPORT_PATH = "../Data/DatePoints/Json/"

def main():
    n, point_array = index_points()

    if EXPORT:
        export_to_csv(point_array)
        return

    Ns, N2s, Nt, N2t, X = calcSpaceTimeCluster(point_array)

    print("\nCounts:")
    print("Ns: " + str(Ns))
    print("N2s: " + str(N2s))
    print("Nt: " + str(Nt))
    print("N2t: " + str(N2t))
    print("X: " + str(X))
    print("n: " + str(n))

    N, E, V = calc_statistics(Ns, N2s, Nt, N2t, X, n)
    std = math.sqrt(V)

    print("\nStatistics:")
    print("N: " + str(N))
    print("E: " + str(E))
    print("V: " + str(V))
    print("Std: " + str(std))

    print("\nZ-score:")
    diff = abs(X - E)
    Z = diff/std
    print("Z: " + str(Z))


def dateFieldValid():
    fields = arcpy.ListFields(POINT_DATA)
    for field in fields:
        if field.name == DATE_FIELD_NAME and field.type == 'Date':
            print("Valid date-field.")
            return True

    print("Date field does not exist. Aborting..")
    return False

def calcSpaceTimeCluster(point_array):
    print("\n\n====================================================")
    print("Init space time cluster calculation")

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

    startTime = time.time()
    for i in range(len(point_array)):
        if i % 100 == 0 and i != 0:
            print(str(i) + " features complete")
            print("Time elapsed: " + str(round(time.time() - startTime, 2)) + " s")

        for j in range(len(point_array)):
            # Do not count i==j matches
            if point_array[i][0] == point_array[j][0]:
                continue

            d_match = False
            t_match = False

            # Check if below distance threshold
            if distance_diff(point_array[i][2], point_array[j][2]) <= D_MAX:
                d_match = True
                Ns += 1

            # Check if below temporal threshold
            if time_diff(point_array[i][1], point_array[j][1]) <= T_MAX:
                t_match = True
                Nt += 1

            if d_match and t_match:
                X += 1

            # Second order counting match j on k
            for k in range(len(point_array)):
                if point_array[i][0] == point_array[k][0] or point_array[j][0] == point_array[k][0]:
                    continue

                if d_match:
                    if distance_diff(point_array[j][2], point_array[k][2]) <= D_MAX:
                        N2s += 1

                if t_match:
                    if time_diff(point_array[j][1], point_array[k][1]) <= T_MAX:
                        N2t += 1

    print("====================================================")
    # Normalize for double counting and return
    return normalize_double_count(Ns, N2s, Nt, N2t, X)


def index_points():
    print("====================================================")
    if not dateFieldValid():
        return

    print("Init point indexing")
    fields = ["FID", DATE_FIELD_NAME, "SHAPE@"]
    # Num points
    n = 0

    point_array = []
    with arcpy.da.SearchCursor(POINT_DATA, fields) as cursor:
        for point in cursor:
            n += 1

            coords = [point[2][0].X, point[2][0].Y]

            # FID, date, [x, y]
            point_attributes = [point[0], point[1], coords]
            point_array.append(point_attributes)

    print("Point indexing finished. \nNum points: " + str(n))
    return n, point_array

def distance_diff(point1_coords, point2_coords):
    x1 = point1_coords[0]
    y1 = point1_coords[1]
    x2 = point2_coords[0]
    y2 = point2_coords[1]

    dist = math.sqrt(math.pow((x1 - x2), 2) + math.pow((y1 - y2), 2))

    return dist


def time_diff(date1, date2):
    delta = date2 - date1
    return abs(delta.days)


def normalize_double_count(Ns, N2s, Nt, N2t, X):
    return Ns/2, N2s/2, Nt/2, N2t/2, X/2


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

def export_to_csv(point_array):
    f_name = EXPORT_PATH + POINT_DATA[:-4] + ".json"
    with open(f_name, 'w') as exportfile:
        json.dump(point_array, exportfile, indent=4, default=str)
        print("\n====================================================")
        print("Exported all points as json.")
        print("New file - " + f_name)
        print("====================================================")


main()



