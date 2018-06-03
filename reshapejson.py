import arcpy
import json
import os


# ___Config___

# Shapefile path
DATAPATH = "../Data/DatePoints/"
arcpy.env.workspace = DATAPATH

# Shapefile name
POINT_DATA = "fiveyears.shp"

# Name of date attribute field
DATE_FIELD_NAME = "Date"


# Path to where the json will be saved
EXPORT_PATH = "../Data/DatePoints/Json/"
OVERWRITE = False

def main():
    n, point_array = index_points()
    export_to_json(n, point_array)


def dateFieldValid():
    fields = arcpy.ListFields(POINT_DATA)
    for field in fields:
        if field.name == DATE_FIELD_NAME and field.type == 'Date':
            print("Valid date-field.")
            return True

    print("Date field does not exist. Aborting..")
    return False


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

    print("Point indexing finished. \nNum points: " + str(n) + "\n")
    return n, point_array


def export_to_json(n, point_array):
    f_name = EXPORT_PATH + POINT_DATA[:-4] + ".json"

    if os.path.isfile(f_name):
        if not OVERWRITE:
            print(f_name + " - already exists, aborting...")
            print("To enable overwriting set OVERWRITE = True")
            return
        else:
            print(f_name + " - already exists, overwriting...")
    else:
        print("Creating new file - " + f_name)

    with open(f_name, 'w') as exportfile:
        json.dump(point_array, exportfile, indent=4, default=str)
        print("\n====================================================")
        print("Exported all " + str(n) + " points as json.")
        print("File - " + f_name)
        print("====================================================")


main()
